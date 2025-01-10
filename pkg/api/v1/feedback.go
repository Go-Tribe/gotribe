// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateFeedBackRequest 指定了 `POST /v1/feedBacks` 接口的请求参数.
type CreateFeedBackRequest struct {
	Title   string `json:"title" valid:"required,stringlength(1|256)"`
	Content string `json:"content" valid:"required,stringlength(1|10240)"`
	Phone   string `json:"phone"`
}

// CreateFeedBackResponse 指定了 `POST /v1/feedBacks` 接口的返回参数.
type CreateFeedBackResponse struct {
	FeedBackID uint `json:"feedBackID"`
}
