// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useraiapikey

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"

	"llmops/pkg/log"
)

// Authenticate validates the Higress Authorization bearer token.
func (u *UserAIAPIKeyController) HigressAuthenticate(c *gin.Context) {
	log.L(c).Info("UserAIAPIKeyController.Authenticate function called.")

	if err := u.srv.UserAIAPIKey().Authenticate(c, authorizationHeader(c)); err != nil {
		log.L(c).Errorw("UserAIAPIKeyController.Authenticate failed", "error", err)
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, gin.H{"status": "ok"})
}

func authorizationHeader(c *gin.Context) string {
	return c.GetHeader("Authorization")
}
