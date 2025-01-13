// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

// WechatTokenResponse 定义了获取access_token的响应结构
type WechatTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// GetWechatAccessToken 获取微信access_token的方法
func GetWechatAccessToken(appid, secret string) (*WechatTokenResponse, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	url := "https://api.weixin.qq.com/cgi-bin/token"
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"grant_type": "client_credential",
			"appid":      appid,
			"secret":     secret,
		}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode(), resp.String())
	}

	var tokenResponse WechatTokenResponse
	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if tokenResponse.ErrCode != 0 {
		return &tokenResponse, fmt.Errorf("error getting access_token: %s", tokenResponse.ErrMsg)
	}

	return &tokenResponse, nil
}
