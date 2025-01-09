// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/comment"
	"gotribe/internal/app/store"
	mw "gotribe/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 注册comment路由.
func CommentRoutes(g *gin.RouterGroup) gin.IRoutes {
	pc := comment.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 comments 路由分组
		commentv1 := v1.Group("/comments", mw.Authn())
		{
			commentv1.POST("", pc.Create)          // 创建评论
			commentv1.POST("reply", pc.Reply)      // 回复评论
			commentv1.GET(":commentID", pc.Get)    // 获取评论详情
			commentv1.PUT(":commentID", pc.Update) // 更新评论
			commentv1.GET("", pc.List)             // 获取评论列表
		}
	}
	return nil
}
