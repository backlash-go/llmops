// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useraiapikey

import (
	"github.com/gin-gonic/gin"

	"llmops/internal/apiserver/deps"
	consolesrv "llmops/internal/apiserver/service/v1/console"
	"llmops/internal/pkg/middleware"
)

// UserAIAPIKeyController creates a user AI API key handler.
type UserAIAPIKeyController struct {
	srv consolesrv.Service
}

// NewUserAIAPIKeyController creates a user AI API key handler.
func NewUserAIAPIKeyController(depsIns *deps.Dependencies) *UserAIAPIKeyController {
	return &UserAIAPIKeyController{
		srv: consolesrv.NewService(depsIns),
	}
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint64)

	return userID, ok
}
