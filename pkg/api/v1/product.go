// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetProductResponse 指定了 `GET /v1/products/{productID}` 接口的返回参数.
type GetProductResponse ProductInfo

// ProductInfo 详细信息.
type ProductInfo struct {
	ProductID     string            `json:"productID"`
	Title         string            `json:"title"`
	ProductNumber string            `json:"productNumber"`
	ProjectID     string            `json:"projectID"`
	Description   string            `json:"description"`
	Image         []string          `json:"image"`
	Video         string            `json:"video"`
	BuyLimit      uint              `json:"buyLimit"`
	CategoryID    string            `json:"categoryID"`
	Content       string            `json:"content"`
	HtmlContent   string            `json:"Htmlcontent"`
	Tags          []*TagInfo        `json:"tags"`
	Skus          []*ProductSKUInfo `json:"skus"`
	CreatedAt     string            `json:"createdAt"`
	UpdatedAt     string            `json:"updatedAt"`
}

// ListProductRequest 指定了 `GET /v1/products` 接口的请求参数.
type ListProductRequest struct {
	CategoryID string `form:"categoryID"  valid:"required"`
	Offset     int    `form:"offset"`
	Limit      int    `form:"limit"`
}

// ListProductResponse 指定了 `GET /v1/products` 接口的返回参数.
type ListProductResponse struct {
	TotalCount int64          `json:"totalCount"`
	Products   []*ProductInfo `json:"products"`
}

// ProductSKUInfo 指定了sku详细信息.
type ProductSKUInfo struct {
	SkuID         string `json:"skuID"`
	Title         string `json:"title"`
	ProjectID     string `json:"projectID"`
	ProductID     string `json:"productID"`
	Image         string `son:"image"`
	Video         string `json:"video"`
	UnitPrice     int    `json:"unitPrice"`
	MarketPrice   int    `json:"marketPrice"`
	Quantity      uint   `json:"quantity"`
	UnitPoint     int    `json:"unitPoint"`
	EnableDefault uint   ` json:"enableDefault"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
