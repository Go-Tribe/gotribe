// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gotribe/pkg/auth"
	"time"

	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

// UserM 是数据库中 user 记录 struct 格式的映射.
type UserM struct {
	gorm.Model
	UserID    string     `gorm:"type:char(10);not null;uniqueIndex;example:字符ID，分布式 ID;" json:"user_id"`
	Username  string     `gorm:"type:varchar(30);not null;uniqueIndex;example:用户名" json:"username"`
	ProjectID string     `gorm:"type:char(10);not null;index;example:项目ID;" json:"project_id"`
	Password  string     `gorm:"type:varchar(255);not null;example:密码" json:"-"`
	Nickname  string     `gorm:"type:varchar(30);not null;example:昵称" json:"nickname"`
	Email     string     `gorm:"type:varchar(30);not null;uniqueIndex;example:邮箱" json:"email"`
	Phone     string     `gorm:"type:varchar(21);not null;uniqueIndex;example:电话" json:"phone"`
	Sex       string     `gorm:"type:char(1);not null;default:M;example:M:男 F:女" json:"sex"`
	Point     float64    `gorm:"-" json:"point"`
	Status    uint8      `gorm:"type:tinyint(1);not null;default:1;example:用户状态，1-正常；2-禁用" json:"status"`
	Birthday  *time.Time `gorm:"type:date;example:'用户生日，格式为YYYY-MM-DD'" json:"birthday"`
	AvatarURL string     `gorm:"type:varchar(255);example:头像地址" json:"avatar_url"`
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
