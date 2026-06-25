package router

import (
	"github.com/gin-gonic/gin"

	userv1 "llmops/internal/apiserver/controller/v1/user"
	"llmops/internal/apiserver/store/mysql"
)

// RegisterLoginRoutes registers unauthenticated login routes on the user controller.
func RegisterLoginRoutes(store mysql.Factory, g *gin.Engine) {
	userController := userv1.NewUserController(store)

	g.GET(userv1.GenericOAuthPath, userController.OauthLogin)
}
