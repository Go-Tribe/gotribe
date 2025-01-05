// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// LoginRequest 指定了 `POST /login` 接口的请求参数.
type LoginRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(6|18)"`
	Password string `json:"password" valid:"required,stringlength(6|255)"`
}

// LoginResponse 指定了 `POST /login` 接口的返回参数.
type LoginResponse struct {
	Token string `json:"token"`
}

// ChangePasswordRequest 指定了 `POST /v1/users/{name}/change-password` 接口的请求参数.
type ChangePasswordRequest struct {
	// 旧密码.
	OldPassword string `json:"oldPassword" valid:"required,stringlength(6|255)"`

	// 新密码.
	NewPassword string `json:"newPassword" valid:"required,stringlength(6|255)"`
}

// CreateUserRequest 指定了 `POST /v1/users` 接口的请求参数.
type CreateUserRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(6|18)"`
	Password string `json:"password" valid:"required,stringlength(6|255)"`
	Nickname string `json:"nickname" valid:"required,stringlength(2|30)"`
	Email    string `json:"email" valid:"required,email"`
	Phone    string `json:"phone"`
}

// GetUserResponse 指定了 `GET /v1/users/{name}` 接口的返回参数.
type GetUserResponse UserInfo

// UserInfo 指定了用户的详细信息.
type UserInfo struct {
	UserID    string `json:"userID"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PostCount int64  `json:"postCount"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ListUserRequest 指定了 `GET /v1/users` 接口的请求参数.
type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// ListUserResponse 指定了 `GET /v1/users` 接口的返回参数.
type ListUserResponse struct {
	TotalCount int64       `json:"totalCount"`
	Users      []*UserInfo `json:"users"`
}

// UpdateUserRequest 指定了 `PUT /v1/users/{name}` 接口的请求参数.
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" valid:"stringlength(2|30)"`
	Email    *string `json:"email" valid:"email"`
	Phone    *string `json:"phone" valid:"stringlength(11|11)"`
}

type UserWhere struct {
	UserID    string
	Username  string
	ProjectID string
}

type WechatMiniLoginRequest struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

type GetUserWxPhoneRequest struct {
	Code string `json:"code"`
}

type AccountWhere struct {
	OpenID string `json:"openID"`
}
