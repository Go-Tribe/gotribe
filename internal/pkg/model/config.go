// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// ConfigM 是数据库中config 记录 struct 格式的映射.
type ConfigM struct {
	gorm.Model
	ConfigID    string `gorm:"type:char(10);not null;uniqueIndex;example:字符ID，分布式 ID;" json:"configID"`
	Alias       string `gorm:"type:varchar(20);uniqueIndex;example:别名" json:"alias"`
	Title       string `gorm:"type:varchar(30);example:标题" json:"title"`
	Description string `gorm:"type:varchar(300);example:描述" json:"description"`
	Type        uint8  `gorm:"type:tinyint;default:1;example:类型，1表示普通配置" json:"type"`
	Info        string `gorm:"type:longtext;example:内容" json:"info"`
	Status      uint8  `gorm:"type:tinyint;example:状态，1-正常；2-禁用" json:"status"`
}

// TableName 用来指定映射的 MySQL 表名.
func (m *ConfigM) TableName() string {
	return "config"
}

// BeforeCreate 在创建数据库记录之前生成ConfigID.
func (m *ConfigM) BeforeCreate(tx *gorm.DB) error {
	m.ConfigID = gid.GenShortID()

	return nil
}
