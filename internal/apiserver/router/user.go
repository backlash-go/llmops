// Package router provides HTTP routing.
package router

import (
	"github.com/gin-gonic/gin"

	userv1 "llmops/internal/apiserver/controller/v1/user"
	"llmops/internal/apiserver/store"
)

func RegisterUserRoutes(store store.Factory, v *gin.RouterGroup) {
	userController := userv1.NewUserController(store)

	users := v.Group("/users")

	users.POST("", userController.Create)

}
