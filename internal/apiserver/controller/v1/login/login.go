// Package login implements authentication handlers.
package login

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	"llmops/pkg/log"
)

const (
	oauthStateCookieName = "oauth_state"
	oauthStateBytes      = 32
)

// Controller handles login requests.
type Controller struct {
	config Config
}

// Config contains the settings required by the generic OAuth login flow.
type Config struct {
	AuthorizationEndpoint string
	ClientID              string
	RedirectURI           string
	Scopes                []string
	StateTTL              time.Duration
	CookieSecure          bool
}

// NewController creates a login controller.
func NewController(config Config) *Controller {
	return &Controller{config: config}
}

// GenericOAuth starts the Keycloak authorization flow and handles its callback.
func (l *Controller) GenericOAuth(c *gin.Context) {
	queryState := c.Query("state")
	authorizationCode := c.Query("code")
	oauthError := c.Query("error")

	if queryState == "" && authorizationCode == "" && oauthError == "" {
		l.startAuthorization(c)

		return
	}

	l.handleCallback(c, queryState, authorizationCode, oauthError)
}

func (l *Controller) startAuthorization(c *gin.Context) {
	if err := l.validateConfiguration(); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, err.Error()), nil)

		return
	}

	state, err := generateState()
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, "generate OAuth state failed"), nil)

		return
	}

	maxAge := int(l.config.StateTTL.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		oauthStateCookieName,
		state,
		maxAge,
		"/ops/login/generic_oauth",
		"",
		l.config.CookieSecure,
		true,
	)

	authorizationURL, err := l.authorizationURL(state)
	if err != nil {
		deleteStateCookie(c, l.config.CookieSecure)
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, err.Error()), nil)

		return
	}

	log.L(c).Info("redirecting user to OAuth authorization endpoint")
	c.Redirect(http.StatusFound, authorizationURL)
}

func (l *Controller) handleCallback(c *gin.Context, queryState, authorizationCode, oauthError string) {
	cookieState, err := c.Cookie(oauthStateCookieName)
	if err != nil ||
		queryState == "" ||
		subtle.ConstantTimeCompare([]byte(queryState), []byte(cookieState)) != 1 {
		deleteStateCookie(c, l.config.CookieSecure)
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "invalid or expired OAuth state"), nil)

		return
	}

	// A state value is single-use, including callbacks that contain an OAuth error.
	deleteStateCookie(c, l.config.CookieSecure)

	if oauthError != "" {
		description := c.Query("error_description")
		log.L(c).Warnf("OAuth authorization failed: error=%s", oauthError)
		core.WriteResponse(
			c,
			errors.WithCode(code.ErrValidation, "OAuth authorization failed: %s", description),
			nil,
		)

		return
	}

	if authorizationCode == "" {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "OAuth authorization code is missing"), nil)

		return
	}

	// The next step is to exchange this one-time code for tokens on the server.
	core.WriteResponse(c, nil, map[string]string{
		"code":   authorizationCode,
		"status": "state_validated",
	})
}

func (l *Controller) authorizationURL(state string) (string, error) {
	endpoint, err := url.Parse(l.config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("parse OAuth authorization endpoint: %w", err)
	}
	if endpoint.Scheme != "https" || endpoint.Host == "" {
		return "", fmt.Errorf("OAuth authorization endpoint must be an absolute HTTPS URL")
	}

	query := endpoint.Query()
	query.Set("client_id", l.config.ClientID)
	query.Set("redirect_uri", l.config.RedirectURI)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(l.config.Scopes, " "))
	query.Set("state", state)
	endpoint.RawQuery = query.Encode()

	return endpoint.String(), nil
}

func (l *Controller) validateConfiguration() error {
	if l.config.AuthorizationEndpoint == "" {
		return fmt.Errorf("OAuth authorization endpoint is not configured")
	}
	if l.config.ClientID == "" {
		return fmt.Errorf("OAuth client ID is not configured")
	}
	if l.config.RedirectURI == "" {
		return fmt.Errorf("OAuth redirect URI is not configured")
	}
	if l.config.StateTTL <= 0 {
		return fmt.Errorf("OAuth state TTL must be greater than zero")
	}

	return nil
}

func generateState() (string, error) {
	randomBytes := make([]byte, oauthStateBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func deleteStateCookie(c *gin.Context, secure bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		oauthStateCookieName,
		"",
		-1,
		"/ops/login/generic_oauth",
		"",
		secure,
		true,
	)
}
