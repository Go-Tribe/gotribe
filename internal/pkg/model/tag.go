// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
)

// TagM 是数据库中 tag 记录 struct 格式的映射.
type TagM struct {
	gorm.Model
	TagID       string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"tagID"`
	Title       string `gorm:"type:varchar(255);uniqueIndex;not null;example:标题" json:"title"`
	Description string `gorm:"not null;size:300;example:描述" json:"description"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *TagM) TableName() string {
	return "tag"
}
