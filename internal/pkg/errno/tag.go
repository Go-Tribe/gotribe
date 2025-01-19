// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrTagNotFound 表示未找到标签信息.
var ErrTagNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.TagNotFound", Message: map[string]string{
	"en": "Tag was not found.",
	"zh": "标签未找到。",
}}
