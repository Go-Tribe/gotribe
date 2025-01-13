// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

// The package model defines the data structures and database models used by the application.
package model

// Import the gorm package to use its ORM features for database operations.
import (
	"gorm.io/gorm"
)

// ProjectM is a struct that maps to the project record in the database.
// It represents the model for the project table, used for database read/write operations.
type ProjectM struct {
	gorm.Model
	// ProjectID is the unique identifier for the project, using a distributed ID generation strategy.
	ProjectID string `gorm:"type:char(10);not null;uniqueIndex:idx_project_project_id;example:字符ID，分布式ID" json:"projectID"`
	// Name is the name of the project, a required field.
	Name string `gorm:"type:varchar(30);not null;example:项目名" json:"name,omitempty"`
	// Title is the title of the website, a required field.
	Title string `gorm:"type:varchar(30);not null;example:网站标题" json:"title,omitempty"`
	// Description provides a brief description of the project.
	Description string `gorm:"type:varchar(300);example:描述" json:"description,omitempty"`
	// Keywords are used for SEO, listing key information of the website.
	Keywords string `gorm:"type:varchar(30);example:网站关键词" json:"keywords,omitempty"`
	// Domain is the domain name of the project.
	Domain string `gorm:"type:varchar(60);example:项目域名" json:"domain,omitempty"`
	// PostURL is the URL pattern for content pages.
	PostURL string `gorm:"type:varchar(300);example:内容链接" json:"postURL,omitempty"`
	// ICP is the record number for the ICP备案.
	ICP string `gorm:"type:varchar(30);example:icp备案信息" json:"icp,omitempty"`
	// PublicSecurity is the record number for the public security备案.
	PublicSecurity string `gorm:"type:varchar(30);example:公安备案" json:"publicSecurity,omitempty"`
	// Author holds the copyright information of the website.
	Author string `gorm:"type:varchar(30);example:网站版权" json:"author,omitempty"`
	// Info contains detailed content or additional information about the project.
	Info string `gorm:"type:longtext;example:内容" json:"info,omitempty"`
	// BaiduAnalytics is the identifier for Baidu Statistics.
	BaiduAnalytics string `gorm:"type:varchar(30);example:百度统计" json:"baiduAnalytics,omitempty"`
	// Favicon is the path to the website's favicon.
	Favicon string `gorm:"type:varchar(255);example:favicon" json:"favicon,omitempty"`
	// NavImage is the path to the navigation image.
	NavImage string `gorm:"type:varchar(255);example:导航图片" json:"navImage,omitempty"`
	// Status indicates the status of the project, where 1 is normal and 2 is disabled.
	Status int8 `gorm:"type:tinyint;not null;default:1;example:状态，1-正常；2-禁用" json:"status,omitempty"`
}

// TableName specifies the name of the database table that the ProjectM struct corresponds to.
// This method overrides the default table naming convention of gorm.
func (m *ProjectM) TableName() string {
	return "project"
}
