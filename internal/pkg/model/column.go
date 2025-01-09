// Copyright 2023 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ColumnM struct {
	gorm.Model
	ColumnID    string `gorm:"type:char(10);not null;uniqueIndex:idx_column_column_id;example:字符ID，分布式ID" json:"columnID"`
	ProjectID   string `gorm:"type:char(10);not null;index;example:项目ID;" json:"projectID"`
	Title       string `gorm:"type:varchar(30);not null;example:标题" json:"title,omitempty"`
	Description string `gorm:"type:varchar(300);example:描述" json:"description,omitempty"`
	Info        string `gorm:"type:longtext;example:内容" json:"info,omitempty"`
	Icon        string `gorm:"type:longtext;example:图标" json:"icon,omitempty"`
	Ext         string `gorm:"type:text;example:'扩展字段'" json:"ext"`
	Status      uint8  `gorm:"type:tinyint;not null;default:1;example:状态，1-正常；2-禁用" json:"status,omitempty"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *ColumnM) TableName() string {
	return "column"
}

// BeforeCreate 在创建数据库记录之前生成ColumnID.
func (p *ColumnM) BeforeCreate(tx *gorm.DB) error {
	p.ColumnID = gid.GenShortID()
	return nil
}
