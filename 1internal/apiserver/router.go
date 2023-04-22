// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"easy-go/1internal/apiserver/controller/v1/user"
	"easy-go/1internal/apiserver/store/mysql"
	"easy-go/1internal/pkg/code"
	"easy-go/1internal/pkg/middleware"
	"easy-go/1internal/pkg/middleware/auth"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}
func printHello(c *gin.Context) {
	core.WriteResponse(c, nil, "hello")
}

func installController(g *gin.Engine) *gin.Engine {
	// Middlewares.
	g.POST("/hello", printHello)

	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	// Refresh time can be longer than token timeout
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)
			userv1.Use(auto.AuthFunc(), middleware.Validation())
			userv1.PUT(":name", userController.Update)
			userv1.GET(":name", userController.Get) // admin api
		}
	}
	return g
}
