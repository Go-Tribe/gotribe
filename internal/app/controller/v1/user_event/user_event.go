// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package userEvent

import (
	"gotribe/internal/app/biz"
	"gotribe/internal/app/store"
)

// UserEventController 是 userEvent 模块在 Controller 层的实现，用来处理用户模块的请求.
type UserEventController struct {
	// a *auth.Authz
	b biz.IBiz
	// pb.UnimplementedGoTribeServer
}

// New 创建一个 userEvent controller.
func New(ds store.IStore) *UserEventController {
	return &UserEventController{b: biz.NewBiz(ds)}
}
