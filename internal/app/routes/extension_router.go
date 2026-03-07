// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/extension"
	"gotribe/internal/app/store"

	"github.com/gin-gonic/gin"
)

// ExtensionRoutes 注册 extension 路由.
func ExtensionRoutes(g *gin.RouterGroup) gin.IRoutes {
	ext := extension.New(store.S)
	v1 := g.Group("/v1")
	{
		extv1 := v1.Group("/extensions")
		{
			extv1.GET("/market", ext.MarketList) // 插件列表，可选 query: categoryID, type, offset, limit
			extv1.GET("/:extensionID", ext.Get)  // 插件详情
		}
	}
	return g
}
