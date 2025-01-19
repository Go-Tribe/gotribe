// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrProjectNotFound 表示未找到项目信息.
var ErrProjectNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.ProjectNotFound", Message: map[string]string{
	"en": "Project was not found.",
	"zh": "项目未找到。",
}}
