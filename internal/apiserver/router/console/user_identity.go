// Package router provides HTTP routing.
package router

import (
	"github.com/gin-gonic/gin"

	useridentityv1 "llmops/internal/apiserver/controller/v1/console/user_identity"
	"llmops/internal/apiserver/deps"
)

// RegisterUserIdentityRoutes registers user identity resource routes.
func RegisterUserIdentityRoutes(depsIns *deps.Dependencies, v *gin.RouterGroup) {
	userIdentityController := useridentityv1.NewUserIdentityController(depsIns)

	userIdentities := v.Group("/user-identities")

	userIdentities.POST("", userIdentityController.Create)
	userIdentities.GET("", userIdentityController.List)
	userIdentities.GET("/:id", userIdentityController.Get)
	userIdentities.PUT("/:id", userIdentityController.Update)
	userIdentities.DELETE("/:id", userIdentityController.Delete)
}
