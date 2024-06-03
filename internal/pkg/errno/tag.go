// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrTagNotFound 表示未找到分类信息.
var ErrTagNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.ErrTagNotFound", Message: "Tag was not found."}
