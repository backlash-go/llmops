// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/apiserver/deps"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/middleware"

	"llmops/internal/apiserver/router"

	// custom gin validators.
	_ "llmops/pkg/validator"
)

func initRouter(g *gin.Engine, depsIns *deps.Dependencies) {
	installMiddleware(g)
	installController(g, depsIns)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine, depsIns *deps.Dependencies) *gin.Engine {
	// Middlewares.

	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	g.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"method": "GET"})
	})

	v1 := g.Group("/ops/api/v1")
	v1.Use(middleware.CookieSession(depsIns.Redis))

	router.RegisterUserRoutes(depsIns, g, v1)

	router.RegisterUserIdentityRoutes(depsIns, v1)

	return g
}
