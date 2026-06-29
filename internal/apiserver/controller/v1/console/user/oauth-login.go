package user

import (
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	"llmops/internal/pkg/session"
	apiv1 "llmops/pkg/api/llmops/v1"
	"llmops/pkg/log"
)

const (
	// GenericOAuthPath is the Keycloak authorization callback path.
	GenericOAuthPath     = "/ops/login/generic_oauth"
	oauthStateCookieName = "llmops_oauth_state"
	oauthStateBytes      = 32
)

// OAuthLoginConfig contains settings required by the Keycloak login flow.
type OAuthLoginConfig struct {
	AuthorizationEndpoint string
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Scopes                []string
	StateTTL              time.Duration
	CookieSecure          bool
}

// OAuthConfig stores the process-wide OAuth login config initialized at startup.
var OAuthConfig OAuthLoginConfig

// OauthLogin starts the Keycloak authorization flow and handles its callback.
func (u *UserController) OauthLogin(c *gin.Context) {
	log.L(c).Info("user OauthLogin function called.")

	queryState := c.Query("state")
	authorizationCode := c.Query("code")
	oauthError := c.Query("error")

	log.L(c).Infof(
		"OAuth request received: has_state=%t has_code=%t has_error=%t",
		queryState != "",
		authorizationCode != "",
		oauthError != "",
	)

	if queryState == "" && authorizationCode == "" && oauthError == "" {
		u.startOAuthAuthorization(c)

		return
	}

	if queryState == "" || authorizationCode == "" || oauthError != "" {
		log.L(c).Warnf(
			"invalid OAuth callback request: has_state=%t has_code=%t error=%s",
			queryState != "",
			authorizationCode != "",
			oauthError,
		)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "invalid OAuth callback request"), nil)

		return
	}

	u.handleOAuthCallback(c, queryState, authorizationCode)
}

func (u *UserController) startOAuthAuthorization(c *gin.Context) {
	config := OAuthConfig
	if err := validateOAuthConfig(config); err != nil {
		log.L(c).Warnf("OAuth config validation failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, err.Error()), nil)

		return
	}

	state, err := generateState()
	if err != nil {
		log.L(c).Errorf("generate OAuth state failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, "generate OAuth state failed"), nil)

		return
	}

	maxAge := int(config.StateTTL.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		oauthStateCookieName,
		state,
		maxAge,
		GenericOAuthPath,
		"",
		config.CookieSecure,
		true,
	)
	log.L(c).Infof("OAuth state cookie set: max_age=%d path=%s secure=%t", maxAge, GenericOAuthPath, config.CookieSecure)

	authorizationURL, err := buildAuthorizationURL(config, state)
	if err != nil {
		log.L(c).Errorf("build OAuth authorization URL failed: %v", err)
		deleteStateCookie(c, config.CookieSecure)
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, err.Error()), nil)

		return
	}

	log.L(c).Info("redirecting user to Keycloak authorization endpoint")
	c.Redirect(http.StatusFound, authorizationURL)
}

func (u *UserController) handleOAuthCallback(c *gin.Context, queryState, authorizationCode string) {
	config := OAuthConfig
	cookieState, err := c.Cookie(oauthStateCookieName)
	log.L(c).Infof(
		"OAuth callback state check: cookie_found=%t query_state_len=%d cookie_state_len=%d",
		err == nil,
		len(queryState),
		len(cookieState),
	)
	if err != nil ||
		queryState == "" ||
		subtle.ConstantTimeCompare([]byte(queryState), []byte(cookieState)) != 1 {
		log.L(c).Warnf("OAuth state validation failed: cookie_error=%v", err)
		deleteStateCookie(c, config.CookieSecure)

		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "invalid or expired OAuth state"), nil)

		return
	}

	deleteStateCookie(c, config.CookieSecure)
	log.L(c).Info("OAuth state validated and cookie deleted")

	if err := validateOAuthIssuer(c.Query("iss"), config.AuthorizationEndpoint); err != nil {
		log.L(c).Warnf("OAuth issuer validation failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), nil)

		return
	}
	log.L(c).Info("OAuth issuer validation passed")

	tokenResponse, err := exchangeOAuthCode(c.Request.Context(), config, authorizationCode)
	if err != nil {
		log.L(c).Warnf("OAuth code exchange exchangeOAuthCode failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), nil)

		return
	}

	log.L(c).Infof(
		"OAuth code exchanged: token_type=%s expires_in=%d has_access_token=%t has_id_token=%t has_refresh_token=%t",
		tokenResponse.TokenType,
		tokenResponse.ExpiresIn,
		tokenResponse.AccessToken != "",
		tokenResponse.IDToken != "",
		tokenResponse.RefreshToken != "",
	)

	claims, err := validateIDToken(c.Request.Context(), config, tokenResponse.IDToken)
	if err != nil {
		log.L(c).Warnf("OAuth id_token validateIDToken failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), nil)

		return
	}

	log.L(c).Info("OAuth id_token validation passed")

	printIDTokenClaims(claims)

	loginRequest, err := oauthLoginRequestFromClaims(claims)
	if err != nil {
		log.L(c).Errorf("build OAuth login request failed: %v", err)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, err.Error()), nil)

		return
	}

	loginResponse, err := u.srv.User().OauthLogin(c.Request.Context(), loginRequest)
	if err != nil {

		log.L(c).Errorw("UserController.OauthLogin failed", "error", err)

		core.WriteResponse(c, err, nil)

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		session.CookieName,
		loginResponse.SessionID,
		loginResponse.SessionMaxAge,
		session.CookiePath,
		"",
		config.CookieSecure,
		true,
	)
	log.L(c).Infof("OAuth session cookie set: max_age=%d path=%s secure=%t", loginResponse.SessionMaxAge, session.CookiePath, config.CookieSecure)

	c.Redirect(http.StatusFound, session.FrontendPath)
}

func oauthLoginRequestFromClaims(claims map[string]interface{}) (*apiv1.OAuthLoginRequest, error) {
	return &apiv1.OAuthLoginRequest{
		Provider:          "keycloak",
		Issuer:            claimString(claims, "iss"),
		Subject:           claimString(claims, "sub"),
		PreferredUsername: claimString(claims, "preferred_username"),
		Email:             claimString(claims, "email"),
		FirstName:         claimString(claims, "given_name"),
		LastName:          claimString(claims, "family_name"),
		DisplayName:       claimString(claims, "name"),
		Avatar:            claimString(claims, "picture"),
		Roles:             claimStringSlice(claims, "roles"),
	}, nil
}

func claimString(claims map[string]interface{}, name string) string {
	value, _ := claims[name].(string)

	return value
}

func claimStringSlice(claims map[string]interface{}, name string) []string {
	value, ok := claims[name]
	if !ok {
		return []string{}
	}

	items, ok := value.([]interface{})
	if !ok {
		stringItems, ok := value.([]string)
		if !ok {
			return []string{}
		}

		return stringItems
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		if value, ok := item.(string); ok {
			values = append(values, value)
		}
	}

	return values
}
