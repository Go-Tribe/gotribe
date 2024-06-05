// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package category

import (
	"gotribe/internal/gotribe/biz"
	"gotribe/internal/gotribe/store"
)

// CategoryController 是 category 模块在 Controller 层的实现，用来处理示例模块的请求.
type CategoryController struct {
	b biz.IBiz
}

// New 创建一个 category controller.
func New(ds store.IStore) *CategoryController {
	return &CategoryController{b: biz.NewBiz(ds)}
}
