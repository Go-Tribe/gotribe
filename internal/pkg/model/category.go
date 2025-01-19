// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
)

// CategoryM 是数据库中 category 记录 struct 格式的映射.
type CategoryM struct {
	gorm.Model
	CategoryID  string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"categoryID"`
	ParentID    *uint  `gorm:"default:0;example:父菜单编号(编号为0时表示根菜单)" json:"parentID"`
	Sort        uint   `gorm:"default:1;example:排序" json:"sort"`
	Icon        string `gorm:"type:varchar(255);example:图标" json:"icon"`
	Title       string `gorm:"type:varchar(255);not null;example:'标题'" json:"title"`
	Path        string `gorm:"type:varchar(100);example:url" json:"path"`
	Hidden      uint   `gorm:"type:tinyint(1);default:1;example:1显示，2隐藏" json:"hidden"`
	Description string `gorm:"type:varchar(300);not null;example:描述" json:"description"`
	Ext         string `gorm:"type:text;example:扩展字段" json:"ext"`
	Status      uint   `gorm:"type:tinyint;not null;default:1;example:状态，1-正常；2-禁用" json:"status,omitempty"`
}

// TableName 用来指定映射的 MySQL 表名.
func (m *CategoryM) TableName() string {
	return "category"
}
