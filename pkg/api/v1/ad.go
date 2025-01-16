// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateAdRequest 指定了 `POST /v1/ads` 接口的请求参数.
type CreateAdRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CreateAdResponse 指定了 `POST /v1/ads` 接口的返回参数.
type CreateAdResponse struct {
	AdID string `json:"adID"`
}

// GetAdResponse 指定了 `GET /v1/ads/{adID}` 接口的返回参数.
type GetAdResponse AdInfo

// UpdateAdRequest 指定了 `PUT /v1/ads` 接口的请求参数.
type UpdateAdRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}

// AdInfo 指定了文章的详细信息.
type AdInfo struct {
	Username    string `json:"username,omitempty"`
	AdID        string `json:"adID,omitempty"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Sort        uint   `json:"sort"`
	URLType     uint   `json:"urlType"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	SceneID     string `json:"sceneID"`
	Ext         string `json:"ext"`
	Image       string `json:"image"`
	Video       string `json:"video"`
	CategoryID  string `json:"categoryID"`
}

// ListAdRequest 指定了 `GET /v1/ads` 接口的请求参数.
type ListAdRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// ListAdResponse 指定了 `GET /v1/ads` 接口的返回参数.
type ListAdResponse struct {
	TotalCount int64     `json:"totalCount"`
	Ads        []*AdInfo `json:"ads"`
}
