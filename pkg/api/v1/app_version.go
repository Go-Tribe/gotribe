// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetLatestReleaseResponse 指定了 `GET /v1/app_version/latest_release` 的返回参数.
// NeedForceUpdate 由服务端根据最新版本的 ForceUpdate 及最低兼容版本号与客户端版本号比较得出.
type GetLatestReleaseResponse struct {
	ProductName             string `json:"productName"`
	Platform                string `json:"platform"`
	VersionCode             int    `json:"versionCode"`
	VersionName             string `json:"versionName"`
	MinSupportedVersionCode int    `json:"minSupportedVersionCode"`
	ForceUpdate             int    `json:"forceUpdate"`
	Title                   string `json:"title"`
	Content                 string `json:"content"`
	DownloadURL             string `json:"downloadUrl"`
	FileSize                string `json:"fileSize"`
	ReleaseDate             string `json:"releaseDate,omitempty"`
	NeedForceUpdate         bool   `json:"needForceUpdate"` // 是否需要强制升级：最新版 ForceUpdate=1 或 客户端版本 < 最新版最低兼容版本
}
