// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/gotribe/controller/v1/example"
	"gotribe/internal/gotribe/store"
	mw "gotribe/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 注册example路由.
func ExampleRoutes(g *gin.RouterGroup) gin.IRoutes {
	pc := example.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 examples 路由分组
		examplev1 := v1.Group("/examples", mw.Authn())
		{
			examplev1.POST("", pc.Create)             // 创建内容
			examplev1.GET(":exampleID", pc.Get)       // 获取内容详情
			examplev1.PUT(":exampleID", pc.Update)    // 更新内容
			examplev1.DELETE("", pc.DeleteCollection) // 批量删除内容
			examplev1.GET("", pc.List)                // 获取内容列表
			examplev1.DELETE(":exampleID", pc.Delete) // 删除内容
		}
	}
	return nil
}
