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

// List lists user identity bindings from the storage.
func (u *UserIdentityController) List(c *gin.Context) {
	log.L(c).Info("user identity list function called.")

	var r apiv1.ListUserIdentityRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	resp, err := u.srv.UserIdentity().List(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
