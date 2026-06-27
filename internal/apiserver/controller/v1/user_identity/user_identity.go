// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useridentity

import (
	"llmops/internal/apiserver/deps"
	srvv1 "llmops/internal/apiserver/service/v1"
)

// UserIdentityController creates a user identity handler used to handle request for user identity resource.
type UserIdentityController struct {
	srv srvv1.Service
}

// NewUserIdentityController creates a user identity handler.
func NewUserIdentityController(depsIns *deps.Dependencies) *UserIdentityController {
	return &UserIdentityController{
		srv: srvv1.NewService(depsIns),
	}
}
