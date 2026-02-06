// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"github.com/gin-gonic/gin"
	userEvent "gotribe/internal/app/controller/v1/user_event"
	"gotribe/internal/app/store"
)

// 注册用户路由.
func UserEventRoutes(g *gin.RouterGroup) gin.IRoutes {
	uc := userEvent.New(store.S)
	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 userEvents 路由分组
		userEventv1 := v1.Group("/report")
		{
			userEventv1.POST("", uc.Create)
		}
	}
	return g
}
