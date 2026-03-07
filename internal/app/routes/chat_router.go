// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/chat"

	"github.com/gin-gonic/gin"
)

// ChatRoutes 注册 chat 路由（LLM 转发）.
func ChatRoutes(g *gin.RouterGroup) gin.IRoutes {
	ctrl := chat.New()
	v1 := g.Group("/v1")
	{
		chatv1 := v1.Group("/chat")
		{
			chatv1.POST("/completions", ctrl.Completions) // 流式(stream=true) 或 非流式(stream=false)
		}
	}
	return g
}
