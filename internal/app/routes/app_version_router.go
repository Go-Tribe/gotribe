// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/app_version"
	"gotribe/internal/app/store"

	"github.com/gin-gonic/gin"
)

// AppVersionRoutes 注册 app_version 路由.
func AppVersionRoutes(g *gin.RouterGroup) gin.IRoutes {
	ctrl := app_version.New(store.S)
	v1 := g.Group("/v1")
	{
		av := v1.Group("/app_version")
		{
			av.GET("/latest_release", ctrl.LatestRelease) // 获取当前产品+平台最新版本，Header: x-product, x-platform, x-platform-version-code, x-platform-version-name
		}
	}
	return g
}
