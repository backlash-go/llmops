// Package router provides HTTP routing.
package router

import (
	"github.com/gin-gonic/gin"

	userv1 "llmops/internal/apiserver/controller/v1/user"
	"llmops/internal/apiserver/store"
)

func RegisterUserRoutes(g *gin.Engine, store store.Factory, v *gin.RouterGroup) {
	userController := userv1.NewUserController(store)
	users := v.Group("/users")

	users.POST("", userController.Create)
	users.GET("", userController.List)
	users.POST("/batch-delete", userController.DeleteCollection)
	users.GET("/:username", userController.Get)
	users.PUT("/:username", userController.Update)
	users.PUT("/:username/change_password", userController.ChangePassword)
}
