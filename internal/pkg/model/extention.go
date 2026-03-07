// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// 插件
type ExtensionM struct {
	gorm.Model
	ExtensionID string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID" json:"extensionID"`
	CategoryID  string `gorm:"type:varchar(10);Index;comment:分类 ID" json:"categoryID"`
	ProjectID   string `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	UserID      string `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Title       string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Description string `gorm:"type:text;comment:描述" json:"description"`
	Ext         string `gorm:"type:text;comment:'扩展字段'" json:"ext"`
	Icon        string `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Tag         string `gorm:"type:varchar(30);comment:tag" json:"tag"`
	Status      uint   `gorm:"type:smallint;not null;default:1;comment:状态，1-可用；2-禁用" json:"status"`
	Source      uint   `gorm:"type:smallint;not null;default:1;comment:来源，1-系统、2-用户" json:"source"`
	PayType     uint   `gorm:"type:smallint;not null;default:1;comment:付费类型，1-免费、2-vip" json:"payType"`
	Type        uint   `gorm:"type:smallint;not null;comment:类型，1-工具、2-skill、3-工作流、4-渠道" json:"type"`
	Url         string `gorm:"type:varchar(300);comment:文档地址" json:"url"`
	ResourceUrl string `gorm:"type:varchar(300);comment:下载地址" json:"resourceUrl"`
}

func (m *ExtensionM) TableName() string {
	return "extension"
}

func (m *ExtensionM) BeforeCreate(tx *gorm.DB) error {
	m.ExtensionID = gid.GenShortID()
	return nil
}
