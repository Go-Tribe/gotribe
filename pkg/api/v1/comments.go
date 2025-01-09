// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateCommentRequest 指定了 `POST /v1/comment` 接口的请求参数.
type CreateCommentRequest struct {
	Content    string `json:"content" valid:"required,stringlength(1|10240)"`
	ObjectID   string `json:"objectID" valid:"required,stringlength(1|10)"`
	ObjectType string `json:"objectType" valid:"required,stringlength(1|10)"`
}

type ReplyCommentRequest struct {
	Content    string `json:"content" valid:"required,stringlength(1|10240)"`
	ObjectID   string `json:"objectID" valid:"required,stringlength(1|10)"`
	ObjectType string `json:"objectType" valid:"required,stringlength(1|10)"`
	ToUserID   string `json:"toUserID" valid:"required,stringlength(1|10)"`
	PID        int    `json:"pid" valid:"required,stringlength(1|10)"`
	ReplyToID  int    `json:"replyToID" valid:"required,stringlength(1|10)"`
}

// CreateCommentResponse 指定了 `POST /v1/comment` 接口的返回参数.
type CreateCommentResponse struct {
	CommentID string `json:"commentID"`
}

// GetCommentResponse 指定了 `GET /v1/comment/{commentID}` 接口的返回参数.
type GetCommentResponse CommentInfo

// UpdateCommentRequest 指定了 `PUT /v1/comment` 接口的请求参数.
type UpdateCommentRequest struct {
	Content *string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CommentInfo 指定了文章的详细信息.
type CommentInfo struct {
	ID          int            `json:"id"`
	CommentID   string         `json:"commentID"`
	Content     string         `json:"content" `
	HtmlContent string         `json:"htmlContent"`
	ObjectID    string         `json:"objectID"`
	ObjectType  string         `json:"objectType"`
	ToUserID    string         `json:"toUserID" `
	PID         int            `json:"pid" `
	ReplyToID   int            `json:"replyToID"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
	UserID      string         `json:"user_id"`
	Nickname    string         `json:"nickname"`
	Avatar      string         `json:"avatar"`
	Replies     []*CommentInfo `json:"replies"`
}

// ListCommentRequest 指定了 `GET /v1/comment` 接口的请求参数.
type ListCommentRequest struct {
	ObjectID string `form:"objectID" valid:"required,stringlength(1|10)"`
	Offset   int    `form:"offset"`
	Limit    int    `form:"limit"`
}

// ListCommentResponse 指定了 `GET /v1/comment` 接口的返回参数.
type ListCommentResponse struct {
	TotalCount int64          `json:"totalCount"`
	Comments   []*CommentInfo `json:"comment"`
}
