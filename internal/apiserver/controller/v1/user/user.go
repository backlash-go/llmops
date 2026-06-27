// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"llmops/internal/apiserver/deps"
	srvv1 "llmops/internal/apiserver/service/v1"
)

// UserController create a user handler used to handle request for user resource.
type UserController struct {
	srv srvv1.Service
}

// NewUserController creates a user handler.
func NewUserController(depsIns *deps.Dependencies) *UserController {
	return &UserController{
		srv: srvv1.NewService(depsIns),
	}
}
