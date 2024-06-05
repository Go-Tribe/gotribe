// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/gotribe/controller/v1/category"
	"gotribe/internal/gotribe/store"

	"github.com/gin-gonic/gin"
)

// 注册category路由.
func CategoryRoutes(g *gin.RouterGroup) gin.IRoutes {
	cf := category.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 categorys 路由分组
		categoryv1 := v1.Group("/categorys")
		{
			categoryv1.GET(":categoryID", cf.Get) // 获取分类详情
		}
	}
	return nil
}
