// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package known

const (
	// 状态1-待支付；2-已支付；3-已发货；4-已收货；5-已取消；6-待退款；7-已退款；8-退款失败
	OrderStatusPendingPayment = 1 // 待支付
	OrderStatusPaid           = 2 // 已支付
	OrderStatusShipped        = 3 // 已发货
	OrderStatusReceived       = 4 // 已收货
	OrderStatusCanceled       = 5 // 已取消
	OrderStatusRefunding      = 6 // 待退款
	OrderStatusRefunded       = 7 // 已退款
	OrderStatusRefundFailed   = 8 // 退款失败

	// 订单类型 订单类型：1-普通订单；2-积分订单
	OrderTypeNormal = 1
	OrderTypePoint  = 2

	// 支付方式 支付方式：1-微信支付；2-支付宝支付；3-积分支付；4-余额支付
	PaymentMethodWechatPay = 1
	PaymentMethodAlipay    = 2
	PaymentMethodPoint     = 3
	PaymentMethodBalance   = 4
)
