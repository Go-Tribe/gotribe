// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package weixin

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"

	"gotribe/pkg/pay"
)

// WeixinPay 实现 pay.Pay，用于获取微信 Native 支付二维码.
type WeixinPay struct {
	client    *core.Client
	appID     string
	mchID     string
	subAppID  string
	subMchID  string
	notifyURL string
}

// New 从配置文件（viper）读取 pay.weixin 后创建微信支付实现（会加载证书，失败返回 error）.
func New(ctx context.Context) (pay.Pay, error) {
	mchID := viper.GetString("pay.weixin.mch-id")
	mchCertSerial := viper.GetString("pay.weixin.mch-cert-serial-number")
	mchAPIv3Key := viper.GetString("pay.weixin.mch-api-v3-key")
	appID := viper.GetString("pay.weixin.app-id")
	subAppID := viper.GetString("pay.weixin.sub-app-id")
	subMchID := viper.GetString("pay.weixin.sub-mch-id")
	notifyURL := viper.GetString("pay.weixin.notify-url")
	privateKeyPath := viper.GetString("pay.weixin.private-key-path")
	publicKeyPath := viper.GetString("pay.weixin.wechat-pay-public-key-path")

	if mchID == "" || appID == "" || notifyURL == "" || privateKeyPath == "" || publicKeyPath == "" ||
		mchCertSerial == "" || mchAPIv3Key == "" {
		return nil, fmt.Errorf("weixin pay options incomplete (check config pay.weixin.*)")
	}
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load merchant private key: %w", err)
	}
	wechatPayPublicKey, err := utils.LoadPublicKeyWithPath(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load wechat pay public key: %w", err)
	}
	authOpt := option.WithWechatPayPublicKeyAuthCipher(
		mchID, mchCertSerial, mchPrivateKey, mchAPIv3Key, wechatPayPublicKey,
	)
	client, err := core.NewClient(ctx, authOpt)
	if err != nil {
		return nil, fmt.Errorf("new wechat pay client: %w", err)
	}
	return &WeixinPay{
		client: client, appID: appID, mchID: mchID,
		subAppID: subAppID, subMchID: subMchID, notifyURL: notifyURL,
	}, nil
}

// GetQRCode 获取支付二维码（微信 Native 下单，返回 code_url）.
func (w *WeixinPay) GetQRCode(ctx context.Context, req *pay.GetQRCodeRequest) (*pay.GetQRCodeResponse, error) {
	if req == nil || req.OutTradeNo == "" || req.Description == "" || req.AmountTotal <= 0 {
		return nil, fmt.Errorf("pay: invalid GetQRCodeRequest")
	}
	svc := native.NativeApiService{Client: w.client}
	resp, _, err := svc.Prepay(ctx, native.PrepayRequest{
		SpAppid:     core.String(w.appID),
		SpMchid:     core.String(w.mchID),
		SubAppid:    core.String(w.subAppID),
		SubMchid:    core.String(w.subMchID),
		Description: core.String(req.Description),
		OutTradeNo:  core.String(req.OutTradeNo),
		NotifyUrl:   core.String(w.notifyURL),
		Amount:      &native.Amount{Total: core.Int64(req.AmountTotal)},
	})
	if err != nil {
		return nil, fmt.Errorf("weixin native prepay: %w", err)
	}
	if resp == nil || resp.CodeUrl == nil {
		return nil, fmt.Errorf("weixin native prepay: empty code_url")
	}
	return &pay.GetQRCodeResponse{CodeURL: *resp.CodeUrl}, nil
}
