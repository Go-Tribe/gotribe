// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/feedback"
	"gotribe/internal/app/store"
	mw "gotribe/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 注册feedback路由.
func FeedbackRoutes(g *gin.RouterGroup) gin.IRoutes {
	pc := feedback.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 feedbacks 路由分组
		feedbackv1 := v1.Group("/feedback", mw.Authn())
		{
			feedbackv1.POST("", pc.Create) // 创建内容
		}
	}
	return g
}
