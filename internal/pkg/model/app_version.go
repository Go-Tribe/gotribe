// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"time"

	"gorm.io/gorm"
)

// AppVersion 应用版本管理表
type AppVersionM struct {
	gorm.Model
	ProductName             string     `gorm:"type:varchar(20);not null;default:'';comment:产品名" json:"productName"`
	Platform                string     `gorm:"type:varchar(20);not null;default:'';comment:平台：ios/android/harmony" json:"platform"`
	VersionCode             int        `gorm:"type:integer;not null;default:0;comment:版本号（整数，用于比较大小，如32）" json:"versionCode"`
	VersionName             string     `gorm:"type:varchar(50);not null;default:'';comment:版本名称（如3.2.0）" json:"versionName"`
	MinSupportedVersionCode int        `gorm:"type:integer;not null;comment:最低兼容版本（低于此必须升级）" json:"minSupportedVersionCode"`
	ForceUpdate             int        `gorm:"type:smallint;default:0;comment:是否强制升级：0=推荐升级，1=强制升级" json:"forceUpdate"`
	Title                   string     `gorm:"type:varchar(255);comment:升级弹窗标题" json:"title"`
	Content                 string     `gorm:"type:text;comment:升级内容描述" json:"content"`
	DownloadURL             string     `gorm:"type:varchar(500);comment:下载地址" json:"downloadUrl"`
	FileSize                string     `gorm:"type:varchar(50);comment:文件大小（如150MB）" json:"fileSize"`
	ReleaseDate             *time.Time `gorm:"type:timestamp;comment:发布时间" json:"releaseDate,omitempty"`
	Status                  uint       `gorm:"type:smallint;not null;default:1;comment:状态，1-有效、2-失效" json:"status"`
}

func (m *AppVersionM) TableName() string {
	return "app_version"
}
