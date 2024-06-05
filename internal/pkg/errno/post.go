// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrPostNotFound 表示未找到文章.
var ErrPostNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PostNotFound", Message: "Post was not found."}
