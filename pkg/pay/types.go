// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package pay

// GetQRCodeRequest 获取支付二维码的通用入参（pay 层定义，各实现方统一使用）.
type GetQRCodeRequest struct {
	Description string `json:"description"` // 商品描述
	OutTradeNo  string `json:"outTradeNo"`  // 商户订单号
	AmountTotal int64  `json:"amountTotal"` // 订单金额，单位：分
	Attach      string `json:"attach"`      // 附加数据，可选
}

// GetQRCodeResponse 获取支付二维码的通用出参（pay 层定义）.
type GetQRCodeResponse struct {
	CodeURL string `json:"codeUrl"` // 二维码链接，用于生成二维码供用户扫码支付
}
