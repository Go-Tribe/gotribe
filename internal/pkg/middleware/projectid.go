// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"

	"github.com/gin-gonic/gin"
)

// ProjectID 是一个 Gin 中间件，用来在每一个 HTTP 请求的 context, response 中注入 `X-Project-ID` 键值对.
func ProjectID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否有 `X-Project-ID`，如果有则复用，没有中断
		projectID := c.Request.Header.Get(known.XPrjectIDKey)

		if projectID == "" {
			core.WriteResponse(c, errno.ErrUnpropjectID, nil)
			c.Abort()
			return
		}

		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(known.XPrjectIDKey, projectID)

		// 将 ProjectID 保存在 HTTP 返回头中，Header 的键为 `X-Project-ID`
		c.Writer.Header().Set(known.XPrjectIDKey, projectID)
		c.Next()
	}
}
