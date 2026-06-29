// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/pkg/code"
	"llmops/internal/pkg/middleware"
	apiv1 "llmops/pkg/api/llmops/v1"
	"llmops/pkg/log"
)

// Get returns the current user profile from the storage.
func (u *UserController) Get(c *gin.Context) {
	log.L(c).Info("UserController.Get function called.")

	userID, ok := currentUserID(c)
	if !ok || userID == 0 {

	    log.L(c).Errorf("UserController.currentUserID  user auth session fail.")


		core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, "user id is missing from session"), nil)

		return
	}

	resp, err := u.srv.User().Get(c, &apiv1.GetUserRequest{UserID: userID})
	if err != nil {

		log.L(c).Errorf("UserController.u.srv.User().Get  fail.",err)


		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint64)

	return userID, ok
}
