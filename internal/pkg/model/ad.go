// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
)

type AdM struct {
	gorm.Model
	AdID        string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"AdID"`
	Title       string `gorm:"type:varchar(255);not null;example:标题" json:"title"`
	Description string `gorm:"not null;size:300;not null;example:描述" json:"description"`
	URL         string `gorm:"type:varchar(255);not null;example:广告链接" json:"url"`
	URLType     uint   `gorm:"type:tinyint;default:1;example:1.链接，2.文章，3.商品" json:"urlType"`
	Sort        uint   `gorm:"type:tinyint;default:1;example:排序" json:"sort"`
	Status      uint   `gorm:"type:tinyint;not null;default:1;example:状态，1-未发布；2-发布" json:"status,omitempty"`
	SceneID     string `gorm:"type:char(10);Index;example:场景 ID" json:"sceneID"`
	Ext         string `gorm:"type:text;example:扩展字段" json:"ext"`
	Image       string `gorm:"type:varchar(255);example:图片地址" json:"image"`
	Video       string `gorm:"type:varchar(255);example:视频地址" json:"video"`
}

func (m *AdM) TableName() string {
	return "ad"
}
