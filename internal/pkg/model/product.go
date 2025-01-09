package model

import "gorm.io/gorm"

type ProductM struct {
	gorm.Model
	ProductID     string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"productID"`
	Title         string `gorm:"type:varchar(255);not null;example:标题" json:"title"`
	ProductNumber string `gorm:"type:varchar(255);not null;example:商品货号" json:"productNumber"`
	ProjectID     string `gorm:"type:varchar(10);Index;example:项目 ID" json:"projectID"`
	Description   string `gorm:"not null;size:300;not null;example:产品卖点/描述" json:"description"`
	Image         string `gorm:"type:varchar(255);not null;example:产品主图" json:"image"`
	Video         string `gorm:"type:varchar(255);not null;example:产品视频" json:"video"`
	BuyLimit      uint   `gorm:"type:tinyint(4);not null;default:1;example:购买限制" json:"buyLimit"`
	CategoryID    string `gorm:"type:char(10);not null;index;example:分类ID" json:"categoryID"`
	ProductSpec   string `gorm:"type:varchar(2048);not null;example:产品规格" json:"productSpec"`
	Content       string `gorm:"type:longtext;example:内容" json:"content"`
	HtmlContent   string `gorm:"type:longtext;example:html内容" json:"Htmlcontent"`
	Tag           string `gorm:"type:varchar(300);not null;example:标签" json:"tag"`
	Enable        uint   `gorm:"type:tinyint(4);not null;default:1;example:是否启用：1-下架；2-上架" json:"enable"`
}

func (p *ProductM) TableName() string {
	return "product"
}
