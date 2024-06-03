// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package project

import (
	"gotribe/internal/gotribe/biz"
	"gotribe/internal/gotribe/store"
)

// ProjectController 是 project 模块在 Controller 层的实现，用来处理示例模块的请求.
type ProjectController struct {
	b biz.IBiz
}

// New 创建一个 project controller.
func New(ds store.IStore) *ProjectController {
	return &ProjectController{b: biz.NewBiz(ds)}
}
