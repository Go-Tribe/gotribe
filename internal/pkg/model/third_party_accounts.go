// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import "gorm.io/gorm"

type ThirdPartyAccountsM struct {
	gorm.Model
	UserID   string `gorm:"type:char(10);Index;example:用户ID" json:"userID"`
	Platform string `gorm:"type:varchar(50);not null;example:平台" json:"platform"`
	BindFlag uint   `gorm:"type:tinyint;default:1;example:是否绑定,2绑定" json:"bindFlag"`
	OpenID   string `gorm:"type:varchar(255);uniqueIndex;not null;example:openID" json:"openID"`
}

func (u *ThirdPartyAccountsM) TableName() string {
	return "third_party_accounts"
}
