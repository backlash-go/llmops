package user

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"llmops/pkg/log"
)

const keycloakAuthPathSuffix = "/protocol/openid-connect/auth"

var oauthHTTPClient = &http.Client{Timeout: 10 * time.Second}

type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

func buildAuthorizationURL(config OAuthLoginConfig, state string) (string, error) {
	endpoint, err := url.Parse(config.AuthorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("parse OAuth authorization endpoint: %w", err)
	}
	if endpoint.Scheme != "https" || endpoint.Host == "" {
		return "", fmt.Errorf("OAuth authorization endpoint must be an absolute HTTPS URL")
	}

	query := endpoint.Query()
	query.Set("client_id", config.ClientID)
	query.Set("redirect_uri", config.RedirectURI)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(config.Scopes, " "))
	query.Set("state", state)
	endpoint.RawQuery = query.Encode()

	return endpoint.String(), nil
}

func exchangeOAuthCode(ctx context.Context, config OAuthLoginConfig, code string) (*oauthTokenResponse, error) {
	tokenEndpoint, err := keycloakOpenIDEndpointFromAuthorizationEndpoint(config.AuthorizationEndpoint, "token")
	if err != nil {
		return nil, err
	}
	log.Infof("OAuth token endpoint resolved: endpoint=%s client_id=%s has_client_secret=%t", tokenEndpoint, config.ClientID, config.ClientSecret != "")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", config.ClientID)
	form.Set("code", code)
	form.Set("redirect_uri", config.RedirectURI)
	if config.ClientSecret != "" {
		form.Set("client_secret", config.ClientSecret)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create OAuth token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange OAuth code: %w", err)
	}
	defer resp.Body.Close()
	log.Infof("OAuth token endpoint responded: status=%d", resp.StatusCode)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))

		return nil, fmt.Errorf("exchange OAuth code failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tokenResponse oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("decode OAuth token response: %w", err)
	}
	if tokenResponse.IDToken == "" {
		return nil, fmt.Errorf("OAuth token response missing id_token")
	}

	return &tokenResponse, nil
}

func validateIDToken(ctx context.Context, config OAuthLoginConfig, rawIDToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(rawIDToken, claims, func(token *jwt.Token) (interface{}, error) {
		log.Infof("OAuth id_token header received: alg=%v kid=%v", token.Header["alg"], token.Header["kid"])
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected OAuth id_token signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("OAuth id_token missing kid header")
		}

		return fetchKeycloakPublicKey(ctx, config.AuthorizationEndpoint, kid)
	}, jwt.WithAudience(config.ClientID))
	if err != nil {
		return nil, fmt.Errorf("parse OAuth id_token: %w", err)
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("OAuth id_token is invalid")
	}

	if err := validateIDTokenClaims(config, claims); err != nil {
		return nil, err
	}
	log.Infof("OAuth id_token claims validated: sub=%v preferred_username=%v email=%v", claims["sub"], claims["preferred_username"], claims["email"])

	return claims, nil
}

func validateIDTokenClaims(config OAuthLoginConfig, claims jwt.MapClaims) error {
	expectedIssuer, err := keycloakIssuerFromAuthorizationEndpoint(config.AuthorizationEndpoint)
	if err != nil {
		return err
	}

	issuer, ok := claims["iss"].(string)
	if !ok || issuer == "" {
		return fmt.Errorf("OAuth id_token missing issuer")
	}
	if subtle.ConstantTimeCompare([]byte(issuer), []byte(expectedIssuer)) != 1 {
		return fmt.Errorf("OAuth id_token issuer mismatch")
	}

	if !claimHasAudience(claims["aud"], config.ClientID) {
		return fmt.Errorf("OAuth id_token audience mismatch")
	}

	if subject, ok := claims["sub"].(string); !ok || subject == "" {
		return fmt.Errorf("OAuth id_token missing subject")
	}

	return nil
}

func claimHasAudience(audience interface{}, clientID string) bool {
	switch value := audience.(type) {
	case string:
		return value == clientID
	case []interface{}:
		for _, item := range value {
			if item == clientID {
				return true
			}
		}
	case []string:
		for _, item := range value {
			if item == clientID {
				return true
			}
		}
	}

	return false
}

func printIDTokenClaims(claims jwt.MapClaims) {
	prettyClaims, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		fmt.Printf("Keycloak id_token claims: %+v\n", claims)

		return
	}

	fmt.Printf("Keycloak id_token claims:\n%s\n", string(prettyClaims))
}

func fetchKeycloakPublicKey(ctx context.Context, authorizationEndpoint, kid string) (*rsa.PublicKey, error) {
	certsEndpoint, err := keycloakOpenIDEndpointFromAuthorizationEndpoint(authorizationEndpoint, "certs")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, certsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create Keycloak JWKS request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Keycloak JWKS: %w", err)
	}
	defer resp.Body.Close()
	log.Infof("Keycloak JWKS endpoint responded: status=%d kid=%s", resp.StatusCode, kid)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))

		return nil, fmt.Errorf("fetch Keycloak JWKS failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode Keycloak JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			log.Infof("Keycloak JWKS matched key: kid=%s alg=%s kty=%s use=%s", key.Kid, key.Alg, key.Kty, key.Use)
			return rsaPublicKeyFromJWK(key)
		}
	}

	return nil, fmt.Errorf("Keycloak JWKS missing kid %s", kid)
}

func rsaPublicKeyFromJWK(key jwkKey) (*rsa.PublicKey, error) {
	if key.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported Keycloak JWK kty: %s", key.Kty)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("decode Keycloak JWK modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("decode Keycloak JWK exponent: %w", err)
	}

	exponent := 0
	for _, value := range eBytes {
		exponent = exponent<<8 + int(value)
	}
	if exponent == 0 {
		return nil, fmt.Errorf("invalid Keycloak JWK exponent")
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: exponent,
	}, nil
}

func validateOAuthIssuer(issuer, authorizationEndpoint string) error {
	if issuer == "" {
		return nil
	}

	expectedIssuer, err := keycloakIssuerFromAuthorizationEndpoint(authorizationEndpoint)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(issuer), []byte(expectedIssuer)) != 1 {
		return fmt.Errorf("invalid OAuth issuer")
	}

	return nil
}

func keycloakIssuerFromAuthorizationEndpoint(authorizationEndpoint string) (string, error) {
	endpoint, err := url.Parse(authorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("parse OAuth authorization endpoint: %w", err)
	}
	if !strings.HasSuffix(endpoint.Path, keycloakAuthPathSuffix) {
		return "", fmt.Errorf("OAuth authorization endpoint is not a Keycloak authorization endpoint")
	}

	endpoint.Path = strings.TrimSuffix(endpoint.Path, keycloakAuthPathSuffix)
	endpoint.RawQuery = ""
	endpoint.Fragment = ""

	return endpoint.String(), nil
}

func keycloakOpenIDEndpointFromAuthorizationEndpoint(authorizationEndpoint, endpointName string) (string, error) {
	endpoint, err := url.Parse(authorizationEndpoint)
	if err != nil {
		return "", fmt.Errorf("parse OAuth authorization endpoint: %w", err)
	}
	if !strings.HasSuffix(endpoint.Path, "/auth") {
		return "", fmt.Errorf("OAuth authorization endpoint is not a Keycloak authorization endpoint")
	}

	endpoint.Path = strings.TrimSuffix(endpoint.Path, "/auth") + "/" + endpointName
	endpoint.RawQuery = ""
	endpoint.Fragment = ""

	return endpoint.String(), nil
}

func validateOAuthConfig(config OAuthLoginConfig) error {
	if config.AuthorizationEndpoint == "" {
		return fmt.Errorf("OAuth authorization endpoint is not configured")
	}
	if config.ClientID == "" {
		return fmt.Errorf("OAuth client ID is not configured")
	}
	if config.RedirectURI == "" {
		return fmt.Errorf("OAuth redirect URI is not configured")
	}
	if config.StateTTL <= 0 {
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
		GenericOAuthPath,
		"",
		secure,
		true,
	)
}
