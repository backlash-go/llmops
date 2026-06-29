// Package router provides HTTP routing.
package router

import (
	"github.com/gin-gonic/gin"

	useraiapikeyv1 "llmops/internal/apiserver/controller/v1/console/user_ai_api_key"
	"llmops/internal/apiserver/deps"
)

// RegisterUserAIAPIKeyHigressRoutes registers user AI API key Higress auth routes.
func RegisterUserAIAPIKeyHigressRoutes(depsIns *deps.Dependencies, v *gin.RouterGroup) {
	userAIAPIKeyController := useraiapikeyv1.NewUserAIAPIKeyController(depsIns)

	v.GET("/higress-auth-key", userAIAPIKeyController.HigressAuthenticate)
}

// RegisterUserAIAPIKeyRoutes registers user AI API key resource routes.
func RegisterUserAIAPIKeyRoutes(depsIns *deps.Dependencies, v *gin.RouterGroup) {
	userAIAPIKeyController := useraiapikeyv1.NewUserAIAPIKeyController(depsIns)

	userAIAPIKeys := v.Group("/user-ai-api-keys")

	userAIAPIKeys.POST("", userAIAPIKeyController.Create)
	userAIAPIKeys.POST("/batch-delete", userAIAPIKeyController.BatchDelete)
}
