package router

import (
	"github.com/gin-gonic/gin"

	loginv1 "llmops/internal/apiserver/controller/v1/login"
	genericoptions "llmops/internal/pkg/options"
)

// RegisterLoginRoutes registers unauthenticated login routes.
func RegisterLoginRoutes(options *genericoptions.OAuthOptions, g *gin.Engine) {
	config := loginv1.Config{}
	if options != nil {
		config = loginv1.Config{
			AuthorizationEndpoint: options.AuthorizationEndpoint,
			ClientID:              options.ClientID,
			RedirectURI:           options.RedirectURI,
			Scopes:                options.Scopes,
			StateTTL:              options.StateTTL,
			CookieSecure:          options.CookieSecure,
		}
	}

	loginController := loginv1.NewController(config)

	g.GET("/ops/login/generic_oauth", loginController.GenericOAuth)
}
