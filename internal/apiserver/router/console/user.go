// Package router provides HTTP routing.
package router

import (
	"github.com/gin-gonic/gin"

	userv1 "llmops/internal/apiserver/controller/v1/console/user"
	"llmops/internal/apiserver/deps"
)

// RegisterUserRoutes registers user login and resource routes.
func RegisterUserRoutes(depsIns *deps.Dependencies, g *gin.Engine, v *gin.RouterGroup) {
	userController := userv1.NewUserController(depsIns)

	g.GET(userv1.GenericOAuthPath, userController.OauthLogin)

	users := v.Group("/users")

	users.GET("/me", userController.Get)
	users.POST("", userController.Create)

}
