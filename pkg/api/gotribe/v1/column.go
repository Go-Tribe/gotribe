// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetColumnResponse 指定了 `GET /v1/columns/{columnID}` 接口的返回参数.
type GetColumnResponse ColumnInfo

// ColumnInfo 指定了配置的详细信息.
type ColumnInfo struct {
	ColumnID    string      `json:"ColumnID"`
	ProjectID   string      `json:"projectID"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Info        string      `json:"info"`
	Icon        string      `json:"icon"`
	Status      uint8       `json:"status"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
	Posts       []*PostInfo `json:"posts"`
	PostCount   int64       `json:"postCount"`
}

// ListColumnRequest 指定了 `GET /v1/columns` 接口的请求参数.
type ListColumnRequest struct {
	Offset    int `form:"offset"`
	Limit     int `form:"limit"`
	PostLimit int `form:"postLimit"`
}

// LisColumnResponse 指定了 `GET /v1/columns` 接口的返回参数.
type LisColumnResponse struct {
	TotalCount int64         `json:"totalCount"`
	Columns    []*ColumnInfo `json:"columns"`
}
