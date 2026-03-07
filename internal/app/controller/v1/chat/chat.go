// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package chat

// ChatController 处理 /v1/chat 相关请求，转发到配置的 LLM 后端.
type ChatController struct{}

// New 创建 chat controller.
func New() *ChatController {
	return &ChatController{}
}
