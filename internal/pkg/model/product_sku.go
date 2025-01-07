package model

import "gorm.io/gorm"

type ProductSKUM struct {
	gorm.Model
	SkuID         string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"skuID"`
	Title         string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	ProjectID     string `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	ProductID     string `gorm:"type:varchar(10);Index;comment:产品ID" json:"productID"`
	Image         string `gorm:"type:varchar(255);not null;comment:产品主图" json:"image"`
	Video         string `gorm:"type:varchar(255);not null;comment:产品视频" json:"video"`
	CostPrice     int    `gorm:"type:int(10);not null;comment:成本价" json:"costPrice"`
	UnitPrice     int    `gorm:"type:int(10);not null;comment:商品价格" json:"unitPrice"`
	MarketPrice   int    `gorm:"type:int(10);not null;comment:市场价格" json:"marketPrice"`
	Quantity      uint   `gorm:"type:int(10);not null;comment:库存" json:"quantity"`
	UnitPoint     int    `gorm:"type:float(20,2);NOT NULL;comment:积分数值" json:"unitPoint"`
	EnableDefault uint   `gorm:"type:tinyint(4);not null;default:1;comment:是否启用：1-正常；2-默认" json:"enableDefault"`
}

func (p *ProductSKUM) TableName() string {
	return "product_sku"
}
