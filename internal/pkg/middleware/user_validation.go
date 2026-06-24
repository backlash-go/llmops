// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"llmops/internal/apiserver/store"
	"llmops/internal/pkg/code"
)

// Validation make sure users have the right resource permission and operation.
func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := isAdmin(c); err != nil {
			switch c.FullPath() {
			case "/api/v1/users":
				if c.Request.Method != http.MethodPost {
					core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					c.Abort()

					return
				}
			case "/api/v1/users/:username", "/api/v1/users/:username/change_password":
				username := c.GetString("username")
				if c.Request.Method == http.MethodDelete ||
					(c.Request.Method != http.MethodDelete && username != c.Param("username")) {
					core.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, ""), nil)
					c.Abort()

					return
				}
			default:
			}
		}

		c.Next()
	}
}

// isAdmin make sure the user is administrator.
// It returns a `github.com/marmotedu/errors.withCode` error.
func isAdmin(c *gin.Context) error {
	username := c.GetString(UsernameKey)
	user, err := store.Client().Users().Get(c, username, metav1.GetOptions{})
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	if user.Status != 1 {
		return errors.WithCode(code.ErrPermissionDenied, "user %s is not active", username)
	}

	return nil
}
