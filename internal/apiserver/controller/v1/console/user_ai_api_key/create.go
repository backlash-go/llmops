// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useraiapikey

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	apiv1 "llmops/pkg/api/llmops/v1"
	"llmops/pkg/log"
)

// Create adds a new user AI API key to the storage.
func (u *UserAIAPIKeyController) Create(c *gin.Context) {
	log.L(c).Info("UserAIAPIKeyController.Create function called.")

	var r apiv1.CreateUserAIAPIKeyRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	userID, ok := currentUserID(c)
	if !ok || userID == 0 {
		log.L(c).Errorf("UserAIAPIKeyController.Create user auth session fail.")
		core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, "user id is missing from session"), nil)

		return
	}
	r.UserID = userID

	resp, err := u.srv.UserAIAPIKey().Create(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
