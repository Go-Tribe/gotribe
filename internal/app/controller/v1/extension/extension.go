// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package extension

import (
	"gotribe/internal/app/biz"
	"gotribe/internal/app/store"
)

// ExtensionController 是 extension 模块在 Controller 层的实现.
type ExtensionController struct {
	b biz.IBiz
}

// New 创建一个 extension controller.
func New(ds store.IStore) *ExtensionController {
	return &ExtensionController{b: biz.NewBiz(ds)}
}
