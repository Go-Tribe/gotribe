// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package weixin

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/spf13/viper"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	client *core.Client
)

// NewClient 从配置文件读取 pay.weixin 并创建微信支付 client，供下单、回调验签等复用.
func NewClient(ctx context.Context) (*core.Client, error) {
	if client != nil {
		return client, nil
	}
	mchID := viper.GetString("pay.weixin.mch-id")
	mchCertSerial := viper.GetString("pay.weixin.mch-cert-serial-number")
	mchAPIv3Key := GetMchAPIv3Key()
	privateKeyPath := viper.GetString("pay.weixin.private-key-path")
	publicKeyPath := viper.GetString("pay.weixin.public-key-path")

	if mchID == "" || mchCertSerial == "" || mchAPIv3Key == "" || privateKeyPath == "" || publicKeyPath == "" {
		return nil, fmt.Errorf("weixin pay config incomplete (pay.weixin.*)")
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
	client, err = core.NewClient(ctx, authOpt)
	if err != nil {
		return nil, fmt.Errorf("new wechat pay client: %w", err)
	}
	return client, nil
}

func GetPublicKeyId() string {
	return viper.GetString("public-key-id")
}

func GetMchAPIv3Key() string {
	return viper.GetString("pay.weixin.mch-api-v3-key")
}

func GetPublicKey() (publicKey *rsa.PublicKey, err error) {
	publicKeyPath := viper.GetString("pay.weixin.public-key-path")
	wechatPayPublicKey, err := utils.LoadPublicKeyWithPath(publicKeyPath)
	return wechatPayPublicKey, err
}
