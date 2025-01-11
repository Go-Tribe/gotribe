// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateOrderRequest 指定了 `POST /v1/orders` 接口的请求参数.
type CreateOrderRequest struct {
	ProductID         string `json:"productID" valid:"required,stringlength(1|10)"`
	ProductSkuID      string `json:"productSkuID" valid:"required,stringlength(1|10)"`
	Quantity          int    `json:"quantity" valid:"required"`
	ConsigneeName     string `json:"consigneeName" valid:"required,stringlength(1|256)"`
	ConsigneePhone    string `json:"consigneePhone" valid:"required,stringlength(1|256)"`
	ConsigneeProvince string `json:"consigneeProvince" valid:"required,stringlength(1|256)"`
	ConsigneeCity     string `json:"consigneeCity" valid:"required,stringlength(1|256)"`
	ConsigneeDistrict string `json:"consigneeDistrict" valid:"required,stringlength(1|256)"`
	ConsigneeStreet   string `json:"consigneeStreet" valid:"required,stringlength(1|256)"`
	ConsigneeAddress  string `json:"consigneeAddress" valid:"required,stringlength(1|256)"`
	Remark            string `json:"remark"`
}

// CreateOrderResponse 指定了 `POST /v1/orders` 接口的返回参数.
type CreateOrderResponse struct {
	OrderID string `json:"orderID"`
}

// GetOrderResponse 指定了 `GET /v1/orders/{orderID}` 接口的返回参数.
type GetOrderResponse OrderInfo

// UpdateOrderRequest 指定了 `PUT /v1/orders` 接口的请求参数.
type UpdateOrderRequest struct {
	Title   *string `json:"title" valid:"stringlength(1|256)"`
	Content *string `json:"content" valid:"stringlength(1|10240)"`
}

// OrderInfo 指定了文章的详细信息.
type OrderInfo struct {
	Username    string `json:"username,omitempty"`
	OrderID     string `json:"order_id,omitempty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ListOrderRequest 指定了 `GET /v1/orders` 接口的请求参数.
type ListOrderRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// ListOrderResponse 指定了 `GET /v1/orders` 接口的返回参数.
type ListOrderResponse struct {
	TotalCount int64        `json:"totalCount"`
	Orders     []*OrderInfo `json:"orders"`
}
