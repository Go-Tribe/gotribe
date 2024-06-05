// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package example

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// Update 更新示例.
func (ctrl *ExampleController) Update(c *gin.Context) {
	log.C(c).Infow("Update example function called")

	var r v1.UpdateExampleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Examples().Update(c, c.GetString(known.XUsernameKey), c.Param("exampleID"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
