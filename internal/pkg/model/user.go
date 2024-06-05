// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gotribe/pkg/auth"

	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// UserM 是数据库中 user 记录 struct 格式的映射.
type UserM struct {
	gorm.Model
	UserID    string `gorm:"type:char(10);not null;uniqueIndex;comment:字符ID，分布式 ID;" json:"userID"`
	ProjectID string `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	Username  string `gorm:"type:varchar(30);not null;uniqueIndex;comment:用户名" json:"username"`
	Password  string `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Nickname  string `gorm:"type:varchar(30);not null;comment:昵称" json:"nickname"`
	Email     string `gorm:"type:varchar(30);not null;comment:邮箱" json:"email"`
	Phone     string `gorm:"type:varchar(21);not null;comment:电话" json:"phone"`
	Sex       string `gorm:"type:char(1);not null;default:M;comment:M:男 F:女" json:"sex"`
	Status    uint8  `gorm:"type:tinyint(1);not null;default:1;comment:用户状态，1-正常；2-禁用" json:"status"`
}

// TableName 用来指定映射的 MySQL 表名.
func (u *UserM) TableName() string {
	return "user"
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.UserID = gid.GenShortID()
	return nil
}
