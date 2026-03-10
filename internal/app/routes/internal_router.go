// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe/internal/app/controller/v1/inner"
	"gotribe/internal/app/store"

	"github.com/gin-gonic/gin"
)

// 内部无鉴权接口
func InternalRoutes(g *gin.RouterGroup) gin.IRoutes {
	inn := inner.New(store.S)
	g.POST("/inner/pay/weixin/notify", inn.WeixinPayNotify)
	return g
}
