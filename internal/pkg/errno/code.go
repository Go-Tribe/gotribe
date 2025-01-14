// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

var (
	// OK 代表请求成功.
	OK = &Errno{HTTP: 200, Code: "", Message: map[string]string{
		"en": "Request successful.",
		"zh": "请求成功",
	}}

	// InternalServerError 表示所有未知的服务器端错误.
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: map[string]string{
		"en": "Internal server error.",
		"zh": "内部服务器错误。",
	}}

	// ErrPageNotFound 表示路由不匹配错误.
	ErrPageNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: map[string]string{
		"en": "Page not found.",
		"zh": "页面未找到。",
	}}

	// ErrBind 表示参数绑定错误.
	ErrBind = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: map[string]string{
		"en": "Error occurred while binding the request body to the struct.",
		"zh": "绑定请求体到结构体时出错。",
	}}

	// ErrInvalidParameter 表示所有验证失败的错误.
	ErrInvalidParameter = &Errno{HTTP: 400, Code: "InvalidParameter", Message: map[string]string{
		"en": "Parameter verification failed.",
		"zh": "参数验证失败。",
	}}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{HTTP: 401, Code: "AuthFailure.SignTokenError", Message: map[string]string{
		"en": "Error occurred while signing the JSON web token.",
		"zh": "签发 JWT Token 时出错。",
	}}

	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{HTTP: 401, Code: "AuthFailure.TokenInvalid", Message: map[string]string{
		"en": "Token was invalid.",
		"zh": "Token 格式错误。",
	}}

	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{HTTP: 401, Code: "AuthFailure.Unauthorized", Message: map[string]string{
		"en": "Unauthorized.",
		"zh": "未授权。",
	}}

	// ErrUnpropjectID 表示请求没有归属项目.
	ErrUnpropjectID = &Errno{HTTP: 401, Code: "ResourceNotFound.UnprojectID", Message: map[string]string{
		"en": "UnprojectID.",
		"zh": "未找到项目ID。",
	}}
)
