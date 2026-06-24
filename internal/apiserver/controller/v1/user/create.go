// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	v1 "llmops/pkg/api/llmops/v1"

	"llmops/internal/pkg/code"
	"llmops/pkg/log"
)

// Create add new user to the storage.
func (u *UserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")

	var r v1.CreateUserRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if r.Username == "" || r.Email == "" {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "username and email are required"), nil)
		return
	}

	// Insert the user to the storage.
	if err := u.srv.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
