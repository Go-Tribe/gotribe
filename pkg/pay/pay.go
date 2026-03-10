// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package pay

import "context"

// Pay 支付能力上层接口，入参出参在 pay 层统一，各支付实现（如 weixin）实现此接口即可扩展.
type Pay interface {
	// GetQRCode 获取支付二维码（Native 扫码付），返回可生成二维码的 URL.
	GetQRCode(ctx context.Context, req *GetQRCodeRequest) (*GetQRCodeResponse, error)
}
