// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetConfigResponse 指定了 `GET /v1/configs/{configID}` 接口的返回参数.
type GetConfigResponse ConfigInfo

// ConfigInfo 指定了配置的详细信息.
type ConfigInfo struct {
	ConfigID    string `json:"config_id"`
	Alias       string `json:"alias"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        uint8  `json:"type"`
	Info        string `json:"info"`
	Status      uint8  `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
