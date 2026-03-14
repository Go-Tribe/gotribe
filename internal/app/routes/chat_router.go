// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/chat"
	"gotribe/internal/app/store"

	"github.com/gin-gonic/gin"
)

// ChatRoutes 注册 chat 路由（LLM 转发 + 对话扣费）.
func ChatRoutes(g *gin.RouterGroup) gin.IRoutes {
	ctrl := chat.New(store.S)
	g.POST("/chat/completions", ctrl.Completions)
	return g
}
