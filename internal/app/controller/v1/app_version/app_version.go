// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package app_version

import (
	"gotribe/internal/app/biz"
	"gotribe/internal/app/store"
)

// AppVersionController 处理应用版本相关请求.
type AppVersionController struct {
	b biz.IBiz
}

// New 创建 app_version controller.
func New(ds store.IStore) *AppVersionController {
	return &AppVersionController{b: biz.NewBiz(ds)}
}
