// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package product

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// Get 获取指定的示例.
func (ctrl *ProductController) Get(c *gin.Context) {
	log.C(c).Infow("Get product function called")

	product, err := ctrl.b.Products().Get(c, c.Param("productID"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, product)
}
