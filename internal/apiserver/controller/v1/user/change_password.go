// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	"llmops/pkg/log"
)

// ChangePasswordRequest defines the ChangePasswordRequest data format.
type ChangePasswordRequest struct {
	// Old password.
	// Required: true
	OldPassword string `json:"oldPassword" binding:"omitempty"`

	// New password.
	// Required: true
	NewPassword string `json:"newPassword" binding:"password"`
}

// ChangePassword change the user's password by the user identifier.
func (u *UserController) ChangePassword(c *gin.Context) {
	log.L(c).Info("change password function called.")

	core.WriteResponse(c, errors.WithCode(code.ErrValidation, "password is not part of the current user model"), nil)
}
