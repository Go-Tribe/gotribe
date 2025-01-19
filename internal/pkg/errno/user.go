// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

var (
	// ErrUserAlreadyExist 代表用户已经存在.
	ErrUserAlreadyExist = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: map[string]string{
		"en": "User already exist.",
		"zh": "用户已存在。",
	}}
	// ErrUserNotFound 表示未找到用户.
	ErrUserNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.UserNotFound", Message: map[string]string{
		"en": "User was not found.",
		"zh": "用户未找到。",
	}}
	// ErrPasswordIncorrect 表示密码不正确.
	ErrPasswordIncorrect = &Errno{HTTP: 401, Code: "InvalidParameter.PasswordIncorrect", Message: map[string]string{
		"en": "Password was incorrect.",
		"zh": "密码不正确。",
	}}

	// ErrPermissionDeny 表示用户权限不足.
	ErrPermissionDeny = &Errno{HTTP: 403, Code: "UserAuthFailure.PermissionDeny", Message: map[string]string{
		"en": "Permission denied.",
		"zh": "权限不足。",
	}}
)
