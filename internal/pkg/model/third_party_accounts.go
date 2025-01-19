// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

// Package model defines the structure and operation logic of third-party account models.
package model

// Import gorm package to use its ORM functionality for database operations.
import "gorm.io/gorm"

// ThirdPartyAccountsM represents the structure of a third-party account model.
// It includes fields for user ID, platform name, bind flag, and third-party unique identifier.
type ThirdPartyAccountsM struct {
	gorm.Model
	UserID   string `gorm:"type:char(10);Index;example:用户ID" json:"userID"`                      // User ID
	Platform string `gorm:"type:varchar(50);not null;example:平台" json:"platform"`                // Platform name
	BindFlag uint   `gorm:"type:tinyint;default:1;example:是否绑定,2绑定" json:"bindFlag"`             // Bind flag (2 indicates bound)
	OpenID   string `gorm:"type:varchar(255);uniqueIndex;not null;example:openID" json:"openID"` // Third-party unique identifier
}

// TableName specifies the table name for the ThirdPartyAccountsM model.
// This maps the model to the "third_party_accounts" table in the database.
func (m *ThirdPartyAccountsM) TableName() string {
	return "third_party_accounts"
}
