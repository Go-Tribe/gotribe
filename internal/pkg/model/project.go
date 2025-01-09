// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gorm.io/gorm"
)

// ProjectM 是数据库中 project 记录 struct 格式的映射.
type ProjectM struct {
	gorm.Model
	ProjectID      string `gorm:"type:char(10);not null;uniqueIndex:idx_project_project_id;example:字符ID，分布式ID" json:"projectID"`
	Name           string `gorm:"type:varchar(30);not null;example:项目名" json:"name,omitempty"`
	Title          string `gorm:"type:varchar(30);not null;example:网站标题" json:"title,omitempty"`
	Description    string `gorm:"type:varchar(300);example:描述" json:"description,omitempty"`
	Keywords       string `gorm:"type:varchar(30);example:网站关键词" json:"keywords,omitempty"`
	Domain         string `gorm:"type:varchar(60);example:项目域名" json:"domain,omitempty"`
	PostURL        string `gorm:"type:varchar(300);example:内容链接" json:"postURL,omitempty"`
	ICP            string `gorm:"type:varchar(30);example:icp备案信息" json:"icp,omitempty"`
	PublicSecurity string `gorm:"type:varchar(30);example:公安备案" json:"publicSecurity,omitempty"`
	Author         string `gorm:"type:varchar(30);example:网站版权" json:"author,omitempty"`
	Info           string `gorm:"type:longtext;example:内容" json:"info,omitempty"`
	BaiduAnalytics string `gorm:"type:varchar(30);example:百度统计" json:"baiduAnalytics,omitempty"`
	Favicon        string `gorm:"type:varchar(255);example:favicon" json:"favicon,omitempty"`
	NavImage       string `gorm:"type:varchar(255);example:导航图片" json:"navImage,omitempty"`
	Status         int8   `gorm:"type:tinyint;not null;default:1;example:状态，1-正常；2-禁用" json:"status,omitempty"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *ProjectM) TableName() string {
	return "project"
}
