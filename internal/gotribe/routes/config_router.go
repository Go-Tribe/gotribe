// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/gotribe/controller/v1/config"
	"gotribe/internal/gotribe/store"

	"github.com/gin-gonic/gin"
)

// 注册config路由.
func ConfigRoutes(g *gin.RouterGroup) gin.IRoutes {
	cf := config.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 configs 路由分组
		configv1 := v1.Group("/configs")
		{
			configv1.GET(":alias", cf.Get) // 获取内容详情
		}
	}
	return nil
}
