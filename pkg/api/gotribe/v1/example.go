// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateExampleRequest 指定了 `POST /v1/examples` 接口的请求参数.
type CreateExampleRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CreateExampleResponse 指定了 `POST /v1/examples` 接口的返回参数.
type CreateExampleResponse struct {
	ExampleID string `json:"exampleID"`
}

// GetExampleResponse 指定了 `GET /v1/examples/{exampleID}` 接口的返回参数.
type GetExampleResponse ExampleInfo

// UpdateExampleRequest 指定了 `PUT /v1/examples` 接口的请求参数.
type UpdateExampleRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}

// ExampleInfo 指定了文章的详细信息.
type ExampleInfo struct {
	Username    string `json:"username,omitempty"`
	ExampleID   string `json:"example_id,omitempty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ListExampleRequest 指定了 `GET /v1/examples` 接口的请求参数.
type ListExampleRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// ListExampleResponse 指定了 `GET /v1/examples` 接口的返回参数.
type ListExampleResponse struct {
	TotalCount int64          `json:"totalCount"`
	Examples   []*ExampleInfo `json:"examples"`
}
