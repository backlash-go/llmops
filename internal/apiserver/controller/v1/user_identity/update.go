// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useridentity

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	apiv1 "llmops/pkg/api/llmops/v1"
	"llmops/pkg/log"
)

// Update updates a user identity binding in the storage.
func (u *UserIdentityController) Update(c *gin.Context) {
	log.L(c).Info("user identity update function called.")

	var r apiv1.UpdateUserIdentityRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}
	if err := c.ShouldBindUri(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}
	if r.ID == 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, "user identity id is required"), nil)

		return
	}

	if err := u.srv.UserIdentity().Update(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
