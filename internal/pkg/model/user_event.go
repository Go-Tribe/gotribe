// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
)

// UserM represents the mapping of the user record in the database as a struct.
type UserEventM struct {
	gorm.Model
	UserID      string `gorm:"type:char(30);not null;comment:用户ID;" json:"userID"`
	ProjectID   string `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	EventType   uint8  `gorm:"type:tinyint(1);not null;index;default:1;comment:事件类型，1-浏览事件；2-点击事件" json:"eventType"`
	EventDetail string `gorm:"type:text;comment:事件详情" json:"eventDetail"`
	Duration    int    `gorm:"type:int(11);comment:事件时长" json:"duration"`
	IP          string `gorm:"type:varchar(255);comment:IP地址" json:"ip"`
	UserAgent   string `gorm:"type:varchar(255);comment:用户代理" json:"userAgent"`
	Referer     string `gorm:"type:varchar(255);comment:来源页面" json:"referer"`
	Platform    string `gorm:"type:varchar(255);comment:平台" json:"platform"`
}

// TableName specifies the MySQL table name to which the struct maps.
func (m *UserEventM) TableName() string {
	return "user_event"
}
