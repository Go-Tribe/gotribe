// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package chat

import (
	"gotribe/internal/app/biz"
	"gotribe/internal/app/store"
)

// ChatController 处理 /v1/chat 相关请求，转发到配置的 LLM 后端并做对话扣费.
type ChatController struct {
	b  biz.IBiz
	ds store.IStore
}

// New 创建 chat controller.
func New(ds store.IStore) *ChatController {
	return &ChatController{b: biz.NewBiz(ds), ds: ds}
}
