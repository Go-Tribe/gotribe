// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package inner

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/api/v1"

	"github.com/gin-gonic/gin"
)

// WeixinPayNotify 微信支付成功回调（内部接口：不需鉴权、不需 X-Project-ID）.
// 请求体为 JSON：{"orderNumber":"xxx"}，回调后完成订单支付状态更新，username 传空字符串.
func (ctrl *InnerController) WeixinPayNotify(c *gin.Context) {
	log.C(c).Infow("weixin pay notify called")

	var r v1.PayOrderRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if r.OrderNumber == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("orderNumber required"), nil)
		return
	}

	err := ctrl.b.Orders().Pay(c, r.OrderNumber, "")
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
