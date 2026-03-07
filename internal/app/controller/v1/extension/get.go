// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package extension

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// Get 根据 extensionID 获取插件详情.
func (ctrl *ExtensionController) Get(c *gin.Context) {
	log.C(c).Infow("Get extension function called")

	extensionID := c.Param("extensionID")
	if extensionID == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}
	resp, err := ctrl.b.Extensions().Get(c, extensionID)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)
}
