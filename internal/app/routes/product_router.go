// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/product"
	"gotribe/internal/app/store"

	"github.com/gin-gonic/gin"
)

// 注册product路由.
func ProductRoutes(g *gin.RouterGroup) gin.IRoutes {
	pdt := product.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 products 路由分组
		productv1 := v1.Group("/products")
		{
			productv1.GET(":productID", pdt.Get) // 获取product详情
			productv1.GET("", pdt.List)          // 获取product详情
		}
	}
	return g
}
