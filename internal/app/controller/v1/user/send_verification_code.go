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
	"gotribe/pkg/email"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SendVerificationCode 发送验证码（如注册前发邮箱验证码）. 发件配置从 config 的 email 段读取.
func (ctrl *UserController) SendVerificationCode(c *gin.Context) {
	log.C(c).Infow("Send verification code called")

	var r v1.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	opts := emailOptionsFromViper()
	expireMinutes := viper.GetInt("email.expire-minutes")
	if expireMinutes <= 0 {
		expireMinutes = 10
	}

	if err := ctrl.b.Users().SendVerificationCode(c, &r, opts, expireMinutes); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}

func emailOptionsFromViper() *email.Options {
	host := viper.GetString("email.host")
	from := viper.GetString("email.from")
	if host == "" && from == "" {
		return nil
	}
	port := viper.GetInt("email.port")
	if port == 0 {
		port = 587
	}
	return &email.Options{
		Host:     host,
		Port:     port,
		From:     from,
		Password: viper.GetString("email.password"),
		UseTLS:   viper.GetBool("email.use-tls"),
	}
}
