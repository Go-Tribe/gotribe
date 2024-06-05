// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package tag

import (
	"gotribe/internal/gotribe/biz"
	"gotribe/internal/gotribe/store"
)

// TagController 是 tag 模块在 Controller 层的实现，用来处理示例模块的请求.
type TagController struct {
	b biz.IBiz
}

// New 创建一个 tag controller.
func New(ds store.IStore) *TagController {
	return &TagController{b: biz.NewBiz(ds)}
}
