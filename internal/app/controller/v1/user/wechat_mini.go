package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/log"
	v1 "gotribe/pkg/api/v1"
	"gotribe/pkg/wechat"
)

func (ctrl *UserController) WxMiniLogin(c *gin.Context) {
	log.C(c).Infow("Wxmini Login function called")

	var r v1.WechatMiniLoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	appID := viper.GetString("wechat.mini-app-id")
	appSecret := viper.GetString("wechat.mini-app-secret")
	miniResp, err := wechat.MiniLogin(appID, appSecret, r.Code)
	if err != nil {
		log.C(c).Errorw("Wxmini Login function called", "error", err)
		core.WriteResponse(c, errno.ErrSignToken, nil)
		return
	}
	// 获取登录信息
	log.C(c).Infow("Wxmini Login function called", "resp", miniResp)
	resp, err := ctrl.b.Users().WxMiniLogin(c, &r, miniResp.OpenID)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}
	core.WriteResponse(c, nil, resp)
}

func (ctrl *UserController) GetWxPhone(c *gin.Context) {
	log.C(c).Infow("Get GetWxPhone function called")
	var r v1.GetUserWxPhoneRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}
	appID := viper.GetString("wechat.mini-app-id")
	appSecret := viper.GetString("wechat.mini-app-secret")
	accessTokenRes, err := wechat.GetWechatAccessToken(appID, appSecret)
	if err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}
	rs, err := wechat.GetPhoneNumber(accessTokenRes.AccessToken, r.Code)
	if err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}
	core.WriteResponse(c, nil, rs)
}
