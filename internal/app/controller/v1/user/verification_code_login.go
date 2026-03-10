// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/api/v1"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// VerificationCodeLogin 验证码登录：校验 Target+Code（trigger=login），有账号则登录，无则自动注册再登录. Target 目前仅支持邮箱.
func (ctrl *UserController) VerificationCodeLogin(c *gin.Context) {
	log.C(c).Infow("Verification code login called")

	var r v1.VerificationCodeLoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}
	// 目前仅支持邮箱
	if !govalidator.IsEmail(r.Target) {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("目前仅支持邮箱"), nil)
		return
	}

	resp, err := ctrl.b.Users().VerificationCodeLogin(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)
}
