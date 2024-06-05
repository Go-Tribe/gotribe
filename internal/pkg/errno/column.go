// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

// ErrConfigNotFound 表示未找到专栏信息.
var ErrColumnNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.ColumnNotFound", Message: "Column was not found."}
