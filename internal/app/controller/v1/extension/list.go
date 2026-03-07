// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package extension

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/api/v1"

	"github.com/gin-gonic/gin"
)

// MarketList 返回插件列表，可选参数 categoryID、type，传则参与查询.
func (ctrl *ExtensionController) MarketList(c *gin.Context) {
	log.C(c).Infow("List extension function called")

	var r v1.ListExtensionRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	resp, err := ctrl.b.Extensions().MarketList(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)
}
