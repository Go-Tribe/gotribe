// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// Create 创建一个新的用户.
func (ctrl *UserController) Create(c *gin.Context) {
	log.C(c).Infow("Create user function called")

	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}
	core.WriteResponse(c, nil, nil)
}
