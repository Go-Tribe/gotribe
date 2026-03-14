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
	Token    string `json:"token"`
	Username string `json:"username"`
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

// SendVerificationCodeRequest 发送验证码请求（如注册前发邮箱验证码）.
type SendVerificationCodeRequest struct {
	Email   string `json:"email" valid:"required,email"`
	Trigger string `json:"trigger" valid:"required"` // 场景：register
}

// RegisterRequest 注册请求，需先通过邮箱验证码校验.
type RegisterRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Code     string `json:"code" valid:"required,stringlength(4|8)"` // 邮箱验证码
	Username string `json:"username" valid:"alphanum,required,stringlength(6|18)"`
	Password string `json:"password" valid:"required,stringlength(6|255)"`
	Nickname string `json:"nickname" valid:"required,stringlength(2|30)"`
	Phone    string `json:"phone"`
}

// GetUserResponse 指定了 `GET /v1/users/{name}` 接口的返回参数.
type GetUserResponse UserInfo

// UserInfo 指定了用户的详细信息.
type UserInfo struct {
	UserID     string  `json:"userID"`
	Username   string  `json:"username"`
	Nickname   string  `json:"nickname"`
	Email      string  `json:"email"`
	Sex        string  `json:"sex"`
	Phone      string  `json:"phone"`
	Point      float64 `json:"point"`
	AvatarURL  string  `json:"avatarURL"`
	Birthday   string  `json:"birthday"`
	Background string  `json:"background"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
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
	Nickname   *string `json:"nickname" valid:"stringlength(2|30)"`
	Email      *string `json:"email" valid:"email"`
	Sex        *string `json:"sex" valid:"stringlength(1|2)"`
	AvatarURL  *string `json:"avatarURL"`
	Birthday   *string `json:"birthday"`
	Phone      *string `json:"phone" valid:"stringlength(11|11)"`
	Ext        *string `json:"ext"`
	Background *string `json:"background"`
}

// VerificationCodeLoginRequest 验证码登录请求（目前仅支持邮箱）.
type VerificationCodeLoginRequest struct {
	Target   string `json:"target" valid:"required"` // 邮箱或手机号，目前仅支持邮箱
	Code     string `json:"code" valid:"required"`
	Key      string `json:"key"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type UserWhere struct {
	UserID    string
	Username  string
	ProjectID string
	Email     string
}

type WechatMiniLoginRequest struct {
	Code string `json:"code"`
}

type GetUserWxPhoneRequest struct {
	Code string `json:"code"`
}

type AccountWhere struct {
	OpenID string `json:"openID"`
}

type GetVerificationCodeLoginResultRequest struct {
	Key string `form:"key"`
}
