// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package core

import (
	"net/http"

	"gotribe/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// ErrResponse 定义了发生错误时的返回消息.
type ErrResponse struct {
	// Code 指定了业务错误码.
	Code string `json:"code"`

	// Message 包含了可以直接对外展示的错误信息.
	Message string `json:"message"`
}

// WriteResponse 写响应到 gin.Context 中。
// lang 参数为可选参数，默认为 "zh"。
func WriteResponse(c *gin.Context, err error, data interface{}, langs ...string) {
	lang := c.GetHeader("Accept - Language")
	// 如果没有设置语言，默认使用英语
	if lang != "en" {
		lang = "zh"
	}
	if len(langs) > 0 {
		lang = langs[0]
	}

	if err != nil {
		hcode, code, message := errno.Decode(err, lang)
		c.JSON(hcode, ErrResponse{
			Code:    code,
			Message: message,
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
