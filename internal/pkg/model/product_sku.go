// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import "gorm.io/gorm"

// ProductSKUM represents the model structure for product SKU information.
// It includes basic information such as SKU ID, title, associated project and product IDs, images, videos, prices, stock,积分 value, and enable status.
type ProductSKUM struct {
	gorm.Model
	// SkuID is the unique identifier for the SKU, using a char type with a maximum length of 10, and is indexed uniquely.
	SkuID string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"skuID"`
	// Title is the title of the SKU, using a varchar type with a maximum length of 255, and is a required field.
	Title string `gorm:"type:varchar(255);not null;example:标题" json:"title"`
	// ProjectID is the identifier of the project to which the SKU belongs, using a varchar type with a maximum length of 10, and is indexed.
	ProjectID string `gorm:"type:varchar(10);Index;example:项目 ID" json:"projectID"`
	// ProductID is the identifier of the product to which the SKU belongs, using a varchar type with a maximum length of 10, and is indexed.
	ProductID string `gorm:"type:varchar(10);Index;example:产品ID" json:"productID"`
	// Image is the main image of the product, using a varchar type with a maximum length of 255, and is a required field.
	Image string `gorm:"type:varchar(255);not null;example:产品主图" json:"image"`
	// Video is the video of the product, using a varchar type with a maximum length of 255, and is a required field.
	Video string `gorm:"type:varchar(255);not null;example:产品视频" json:"video"`
	// CostPrice is the cost price of the product, using an int type, and is a required field.
	CostPrice int `gorm:"type:int(10);not null;example:成本价" json:"costPrice"`
	// UnitPrice is the selling price of the product, using an int type, and is a required field.
	UnitPrice int `gorm:"type:int(10);not null;example:商品价格" json:"unitPrice"`
	// MarketPrice is the market price of the product, using an int type, and is a required field.
	MarketPrice int `gorm:"type:int(10);not null;example:市场价格" json:"marketPrice"`
	// Quantity is the stock quantity of the product, using an uint type, and is a required field.
	Quantity uint `gorm:"type:int(10);not null;example:库存" json:"quantity"`
	// UnitPoint is the积分 value of the product, using a float type, and is a required field.
	UnitPoint int `gorm:"type:float(20,2);NOT NULL;example:积分数值" json:"unitPoint"`
	// EnableDefault indicates the enable status of the SKU, using a tinyint type, with a default value of 1, where 1 is normal and 2 is default.
	EnableDefault uint `gorm:"type:tinyint(4);not null;default:1;example:是否启用：1-正常；2-默认" json:"enableDefault"`
}

// TableName specifies the table name for the ProductSKUM model in the database, returning "product_sku".
func (m *ProductSKUM) TableName() string {
	return "product_sku"
}
