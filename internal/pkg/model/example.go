// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// ExampleM 是数据库中 comment 记录 struct 格式的映射.
type ExampleM struct {
	gorm.Model
	ExampleID   string `gorm:"type:varchar(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"exampleID"`
	Username    string `gorm:"type:varchar(30);not null;index:idx_username;example:用户名" json:"username"`
	Title       string `gorm:"type:varchar(255);not null;example:标题" json:"title"`
	Content     string `gorm:"not null;type:longtext;example:内容" json:"content"`
	Description string `gorm:"not null;size:300;example:描述" json:"description"`
	Status      uint8  `gorm:"type:tinyint(1);not null;default:1;example:用户状态，0-禁用；1-启用" json:"status"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *ExampleM) TableName() string {
	return "comment"
}

// BeforeCreate 在创建数据库记录之前生成 ExampleID.
func (p *ExampleM) BeforeCreate(tx *gorm.DB) error {
	p.ExampleID = gid.GenShortID()

	return nil
}
