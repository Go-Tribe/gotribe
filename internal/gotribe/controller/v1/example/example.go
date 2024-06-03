// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package example

import (
	"gotribe/internal/gotribe/biz"
	"gotribe/internal/gotribe/store"
)

// ExampleController 是 example 模块在 Controller 层的实现，用来处理示例模块的请求.
type ExampleController struct {
	b biz.IBiz
}

// New 创建一个 example controller.
func New(ds store.IStore) *ExampleController {
	return &ExampleController{b: biz.NewBiz(ds)}
}
