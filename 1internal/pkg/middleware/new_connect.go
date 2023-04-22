// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"easy-go/2pkg/log"

	"github.com/gin-gonic/gin"
)

// NewConnect is a middleware that create new connect.
func NewConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("alive -1")
		c.Next()
		log.Debugf("alive 0")
		c.GetString(UsernameKey)
		log.Debugf("alive 1,%v", c.GetString(UsernameKey))
	}
}
