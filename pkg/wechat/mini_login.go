package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

// MiniLoginResponse 定义了 code2Session 接口的响应结构
type MiniLoginResponse struct {
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	OpenID     string `json:"openid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// MiniLogin 登录凭证校验。通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程。
func MiniLogin(appid, secret, code string) (*MiniLoginResponse, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	resp, err := client.R().Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode(), resp.String())
	}

	var miniLoginResponse MiniLoginResponse
	err = json.Unmarshal(resp.Body(), &miniLoginResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if miniLoginResponse.ErrCode != 0 {
		return &miniLoginResponse, fmt.Errorf("error getting session info: %s", miniLoginResponse.ErrMsg)
	}

	return &miniLoginResponse, nil
}
