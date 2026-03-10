// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package inner

import (
	"gotribe/pkg/pay/weixin"
	"net/http"

	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

var (
	wechatPayNotifyHandler *notify.Handler
)

// WeixinPayNotify 微信支付成功回调（内部接口：不需鉴权、不需 X-Project-ID）.
// 使用 pkg/pay/weixin 的 NewNotifyHandler（NewRSANotifyHandler）验签并解密，从通知中取订单号并完成订单支付状态，username 传空字符串.
func (ctrl *InnerController) WeixinPayNotify(c *gin.Context) {
	log.C(c).Infow("weixin pay notify called")

	transaction := new(payments.Transaction)

	handler, err := GetNotifyHandler(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "获取Handler失败"})
		return
	}
	notifyReq, err := handler.ParseNotifyRequest(c.Request.Context(), c.Request, transaction)
	if err != nil {
		log.C(c).Warnw("微信支付回调验证失败", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "FAIL", "message": "验证失败"})
		return
	}

	// 打印通知摘要（可选）
	log.C(c).Infow("收到通知摘要: %s", notifyReq.Summary)
	if transaction.OutTradeNo == nil {
		log.C(c).Warnw("交易数据中缺少OutTradeNo")
		c.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL", "message": "数据解析错误"})
		return
	}

	orderNumber := *transaction.OutTradeNo
	transactionId := ""
	if transaction.TransactionId != nil {
		transactionId = *transaction.TransactionId
	}
	log.C(c).Infow("订单 %s 支付成功，微信支付单号: %s", orderNumber, transactionId)

	err = ctrl.b.Orders().Pay(c, orderNumber, "")
	if err != nil {
		log.C(c).Errorw("order pay failed", "orderNumber", orderNumber, "err", err)
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}

func GetNotifyHandler(c *gin.Context) (*notify.Handler, error) {
	client, err := weixin.NewClient(c)
	if err != nil {
		return nil, err
	}
	_ = client // 如果你后续需要调用其他API，可以使用这个client

	mchAPIv3Key := weixin.GetMchAPIv3Key()
	publicKeyID := weixin.GetPublicKeyId()
	publicKey, err := weixin.GetPublicKey()
	if err != nil {
		return nil, err
	}

	// 创建公钥验签器
	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(publicKeyID, *publicKey)

	// 初始化通知处理器
	wechatPayNotifyHandler, err = notify.NewRSANotifyHandler(mchAPIv3Key, verifier)

	return wechatPayNotifyHandler, nil
}
