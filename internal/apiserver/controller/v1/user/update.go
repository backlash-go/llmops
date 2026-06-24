// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	"llmops/pkg/log"
)

// Update update a user info by the user identifier.
func (u *UserController) Update(c *gin.Context) {
	log.L(c).Info("update user function called.")

	var r model.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	user, err := u.srv.Users().Get(c, c.Param("username"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	if r.Email != "" {
		user.Email = r.Email
	}
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Avatar = r.Avatar
	user.Status = r.Status
	user.LastLoginAt = r.LastLoginAt

	// Save changed fields.
	if err := u.srv.Users().Update(c, user, metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
