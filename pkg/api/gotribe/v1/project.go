// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetProjectResponse 指定了 `GET /v1/project/{projectID}` 接口的返回参数.
type GetProjectResponse ProjectInfo

// ProjectInfo 指定了项目详细信息.
type ProjectInfo struct {
	ProjectID      string `json:"projectID,omitempty"`
	Name           string `json:"name,omitempty"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Keywords       string `json:"keywords,omitempty"`
	Domain         string `json:"domain,omitempty"`
	PostURL        string `json:"postURL,omitempty"`
	ICP            string `json:"icp,omitempty"`
	PublicSecurity string `json:"publicSecurity,omitempty"`
	Author         string `json:"author,omitempty"`
	Info           string `json:"info,omitempty"`
	BaiduAnalytics string `json:"baiduAnalytics,omitempty"`
	Favicon        string `json:"favicon,omitempty"`
	NavImage       string `json:"navImage,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}
