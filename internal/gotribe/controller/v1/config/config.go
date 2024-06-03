// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package config

import (
	"gotribe/internal/gotribe/biz"
	"gotribe/internal/gotribe/store"
)

// ConfigController 是 config 模块在 Controller 层的实现，用来处理示例模块的请求.
type ConfigController struct {
	b biz.IBiz
}

// New 创建一个 config controller.
func New(ds store.IStore) *ConfigController {
	return &ConfigController{b: biz.NewBiz(ds)}
}
