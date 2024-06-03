// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreatePostRequest 指定了 `POST /v1/posts` 接口的请求参数.
type CreatePostRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CreatePostResponse 指定了 `POST /v1/posts` 接口的返回参数.
type CreatePostResponse struct {
	PostID string `json:"postID"`
}

// GetPostResponse 指定了 `GET /v1/posts/{postID}` 接口的返回参数.
type GetPostResponse PostInfo

// UpdatePostRequest 指定了 `PUT /v1/posts` 接口的请求参数.
type UpdatePostRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}

// PostInfo 指定了文章的详细信息.
type PostInfo struct {
	Author      string       `json:"author,omitempty"`
	PostID      string       `json:"postID,omitempty"`
	ColumnID    string       `json:"columnID,omitempty"`
	Icon        string       `json:"icon"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	HtmlContent string       `json:"htmlContent"`
	Tag         string       `json:"tag"`
	Description string       `json:"description"`
	Type        uint         `json:"type"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
	Category    CategoryInfo `json:"category"`
	Tags        []*TagInfo   `json:"tags"`
}

// ListPostRequest 指定了 `GET /v1/posts` 接口的请求参数.
type ListPostRequest struct {
	Offset     int    `form:"offset"`
	Limit      int    `form:"limit"`
	Type       string `form:"type"`
	CategoryID string `form:"categoryID"`
	ColumnID   string `form:"columnID"`
	PostID     string `form:"postID"`
	Query      string `form:"query"`
	TagID      string `form:"tagID"`
	IsTop      int    `form:"isTop" valid:"oneof=0 1"`
}

// ListPostResponse 指定了 `GET /v1/posts` 接口的返回参数.
type ListPostResponse struct {
	TotalCount int64       `json:"totalCount"`
	Posts      []*PostInfo `json:"posts"`
}

type PostQueryParams struct {
	PostID     string
	Author     string
	ProjectID  string
	Type       int
	Status     int
	ColumnID   int
	CategoryID int
}
