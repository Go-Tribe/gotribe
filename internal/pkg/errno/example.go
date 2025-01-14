// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrExampleNotFound 表示未找到示例.
var ErrExampleNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.ExampleNotFound", Message: map[string]string{
	"en": "Example was not found.",
	"zh": "示例未找到。",
}}
