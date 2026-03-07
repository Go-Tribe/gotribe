// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// ListExtensionRequest 指定了 `GET /v1/extensions` 接口的请求参数.
// CategoryID、Type 为可选，传则参与查询.
type ListExtensionRequest struct {
	CategoryID string `form:"categoryID"`
	Type       uint   `form:"type"`
	Offset     int    `form:"offset"`
	Limit      int    `form:"limit"`
}

// ListExtensionResponse 指定了 `GET /v1/extensions` 接口的返回参数.
type ListExtensionResponse struct {
	TotalCount int64            `json:"totalCount"`
	Extensions []*ExtensionInfo `json:"extensions"`
}

// GetExtensionResponse 指定了 `GET /v1/extensions/:extensionID` 接口的返回参数.
type GetExtensionResponse ExtensionInfo

// ExtensionInfo 插件详情/列表项.
type ExtensionInfo struct {
	ExtensionID string `json:"extensionID"`
	CategoryID  string `json:"categoryID"`
	ProjectID   string `json:"projectID"`
	UserID      string `json:"userID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Ext         string `json:"ext"`
	Icon        string `json:"icon"`
	Tag         string `json:"tag"`
	Status      uint   `json:"status"`
	Source      uint   `json:"source"`
	PayType     uint   `json:"payType"`
	Type        uint   `json:"type"`
	Url         string `json:"url"`
	ResourceUrl string `json:"resourceUrl"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
