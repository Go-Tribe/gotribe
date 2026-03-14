// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package weixin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"encoding/json"

	"gotribe/pkg/pay"

	"github.com/spf13/viper"
)

// WeixinPay 实现 pay.Pay，用于获取微信 Native 支付二维码.
type WeixinPay struct {
	appID               string
	mchID               string
	mchCertSerialNumber string
	mchApiV3Key         string
	notifyURL           string
	privateKeyPath      string
	publicKeyPath       string
	publicKeyId         string
}

// New 从配置文件（viper）读取 pay.weixin 后创建微信支付实现（通过 weixin_client.NewClient 初始化 client）.
func New(ctx context.Context) (pay.Pay, error) {
	appID := viper.GetString("pay.weixin.app-id")
	mchID := viper.GetString("pay.weixin.mch-id")
	mchCertSerialNumber := viper.GetString("pay.weixin.mch-cert-serial-number")
	mchApiV3Key := viper.GetString("pay.weixin.mch-api-v3-key")
	notifyURL := viper.GetString("pay.weixin.notify-url")
	privateKeyPath := viper.GetString("pay.weixin.private-key-path")
	publicKeyPath := viper.GetString("pay.weixin.public-key-path")
	publicKeyId := viper.GetString("pay.weixin.public-key-id")
	if appID == "" || mchID == "" || notifyURL == "" || mchCertSerialNumber == "" || mchApiV3Key == "" ||
		privateKeyPath == "" || publicKeyPath == "" || publicKeyId == "" {
		return nil, fmt.Errorf("weixin pay options incomplete (check config pay.weixin.*)")
	}

	return &WeixinPay{
		appID:               appID,
		mchID:               mchID,
		notifyURL:           notifyURL,
		mchCertSerialNumber: mchCertSerialNumber,
		mchApiV3Key:         mchApiV3Key,
		privateKeyPath:      privateKeyPath,
		publicKeyPath:       publicKeyPath,
		publicKeyId:         publicKeyId,
	}, nil
}

// GetQRCode 获取支付二维码（微信 Native 下单，返回 code_url）.
func (w *WeixinPay) GetQRCode(ctx context.Context, req *pay.GetQRCodeRequest) (*pay.GetQRCodeResponse, error) {
	if req == nil || req.OutTradeNo == "" || req.Description == "" || req.AmountTotal <= 0 {
		return nil, fmt.Errorf("pay: invalid GetQRCodeRequest")
	}

	config, err := CreateMchConfig(w.mchID, w.mchCertSerialNumber, w.privateKeyPath, w.publicKeyId, w.publicKeyPath)
	if err != nil {
		return nil, err
	}

	request := &CommonPrepayRequest{
		Appid:       String(w.appID),
		Mchid:       String(w.mchID),
		Description: String(req.Description),
		OutTradeNo:  String(req.OutTradeNo),
		TimeExpire:  Time(time.Now()),
		//Attach:        String("自定义数据说明"),
		NotifyUrl:     String(w.notifyURL),
		GoodsTag:      String("WXG"),
		SupportFapiao: Bool(false),
		Amount: &CommonAmountInfo{
			Total:    Int64(req.AmountTotal),
			Currency: String("CNY"),
		},
	}

	response, err := NativePrepay(config, request)
	if err != nil {
		fmt.Printf("请求失败: %+v\n", err)
		return nil, err
	}

	fmt.Printf("请求成功: %+v\n", *response.CodeUrl)
	return &pay.GetQRCodeResponse{CodeURL: *response.CodeUrl}, nil
}

func NativePrepay(config *MchConfig, request *CommonPrepayRequest) (response *DirectApiv3DirectNativePrepayResponse, err error) {
	const (
		host   = "https://api.mch.weixin.qq.com"
		method = "POST"
		path   = "/v3/pay/transactions/native"
	)

	reqUrl, err := url.Parse(fmt.Sprintf("%s%s", host, path))
	if err != nil {
		return nil, err
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequest(method, reqUrl.String(), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Wechatpay-Serial", config.WechatPayPublicKeyId())
	httpRequest.Header.Set("Content-Type", "application/json")
	authorization, err := BuildAuthorization(config.MchId(), config.CertificateSerialNo(), config.PrivateKey(),
		method, reqUrl.RequestURI(), reqBody)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Authorization", authorization)

	client := &http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	respBody, err := ExtractResponseBody(httpResponse)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		// 2XX 成功，验证应答签名
		err = ValidateResponse(
			config.WechatPayPublicKeyId(),
			config.WechatPayPublicKey(),
			&httpResponse.Header,
			respBody,
		)
		if err != nil {
			return nil, err
		}
		response := &DirectApiv3DirectNativePrepayResponse{}
		if err := json.Unmarshal(respBody, response); err != nil {
			return nil, err
		}

		return response, nil
	} else {
		return nil, NewApiException(
			httpResponse.StatusCode,
			httpResponse.Header,
			respBody,
		)
	}
}

type CommonPrepayRequest struct {
	Appid         *string           `json:"appid,omitempty"`
	Mchid         *string           `json:"mchid,omitempty"`
	Description   *string           `json:"description,omitempty"`
	OutTradeNo    *string           `json:"out_trade_no,omitempty"`
	TimeExpire    *time.Time        `json:"time_expire,omitempty"`
	Attach        *string           `json:"attach,omitempty"`
	NotifyUrl     *string           `json:"notify_url,omitempty"`
	GoodsTag      *string           `json:"goods_tag,omitempty"`
	SupportFapiao *bool             `json:"support_fapiao,omitempty"`
	Amount        *CommonAmountInfo `json:"amount,omitempty"`
	Detail        *CouponInfo       `json:"detail,omitempty"`
	SceneInfo     *CommonSceneInfo  `json:"scene_info,omitempty"`
	SettleInfo    *SettleInfo       `json:"settle_info,omitempty"`
}

type DirectApiv3DirectNativePrepayResponse struct {
	CodeUrl *string `json:"code_url,omitempty"`
}

type CommonAmountInfo struct {
	Total    *int64  `json:"total,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

type CouponInfo struct {
	CostPrice   *int64        `json:"cost_price,omitempty"`
	InvoiceId   *string       `json:"invoice_id,omitempty"`
	GoodsDetail []GoodsDetail `json:"goods_detail,omitempty"`
}

type CommonSceneInfo struct {
	PayerClientIp *string    `json:"payer_client_ip,omitempty"`
	DeviceId      *string    `json:"device_id,omitempty"`
	StoreInfo     *StoreInfo `json:"store_info,omitempty"`
}

type SettleInfo struct {
	ProfitSharing *bool `json:"profit_sharing,omitempty"`
}

type GoodsDetail struct {
	MerchantGoodsId  *string `json:"merchant_goods_id,omitempty"`
	WechatpayGoodsId *string `json:"wechatpay_goods_id,omitempty"`
	GoodsName        *string `json:"goods_name,omitempty"`
	Quantity         *int64  `json:"quantity,omitempty"`
	UnitPrice        *int64  `json:"unit_price,omitempty"`
}

type StoreInfo struct {
	Id       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	AreaCode *string `json:"area_code,omitempty"`
	Address  *string `json:"address,omitempty"`
}
