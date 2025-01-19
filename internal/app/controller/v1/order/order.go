// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package order

import (
	"gotribe/internal/app/biz"
	"gotribe/internal/app/store"
)

// OrderController 是 comment 模块在 Controller 层的实现，用来处理示例模块的请求.
type OrderController struct {
	b biz.IBiz
}

// New 创建一个 comment controller.
func New(ds store.IStore) *OrderController {
	return &OrderController{b: biz.NewBiz(ds)}
}
