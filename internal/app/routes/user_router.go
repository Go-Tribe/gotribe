// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/user"
	"gotribe/internal/app/store"
	mw "gotribe/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 注册用户路由.
func UserRoutes(g *gin.RouterGroup) gin.IRoutes {
	uc := user.New(store.S)
	g.POST("/login", uc.Login)
	g.POST("/wxmini/login", uc.WxMiniLogin)
	g.POST("/wxmini/phone", uc.GetWxPhone)
	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)                             // 创建用户
			userv1.PUT(":name/change-password", uc.ChangePassword) // 修改用户密码
			userv1.Use(mw.Authn())
			userv1.GET(":name", uc.Get)                // 获取用户详情
			userv1.PUT(":name", uc.Update)             // 更新用户
			userv1.GET("", uc.List)                    // 列出用户列表
			userv1.DELETE(":name", uc.Delete)          // 删除用户
			userv1.POST("/upload", uc.UploadResources) //上传图片
		}
	}
	return nil
}
