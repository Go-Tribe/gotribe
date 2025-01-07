// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/ad"
	"gotribe/internal/app/store"
	mw "gotribe/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 注册ad路由.
func AdRoutes(g *gin.RouterGroup) gin.IRoutes {
	pc := ad.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 ads 路由分组
		adv1 := v1.Group("/ads", mw.Authn())
		{
			adv1.GET(":sceneID", pc.List)
		}
	}
	return nil
}
