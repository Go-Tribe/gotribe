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

// UserM represents the mapping of the user record in the database as a struct.
type UserM struct {
	gorm.Model
	UserID    string     `gorm:"type:char(10);not null;uniqueIndex;example:Character ID, distributed ID;" json:"user_id"`
	Username  string     `gorm:"type:varchar(30);not null;uniqueIndex;example:Username" json:"username"`
	ProjectID string     `gorm:"type:char(10);not null;index;example:Project ID;" json:"project_id"`
	Password  string     `gorm:"type:varchar(255);not null;example:Password" json:"-"`
	Nickname  string     `gorm:"type:varchar(30);not null;example:Nickname" json:"nickname"`
	Email     string     `gorm:"type:varchar(30);not null;uniqueIndex;example:Email" json:"email"`
	Phone     string     `gorm:"type:varchar(21);not null;uniqueIndex;example:Phone" json:"phone"`
	Sex       string     `gorm:"type:char(1);not null;default:M;example:M:Male F:Female" json:"sex"`
	Point     *float64   `gorm:"-" json:"point"`
	Status    uint8      `gorm:"type:tinyint(1);not null;default:1;example:User status, 1-Normal; 2-Disabled" json:"status"`
	Birthday  *time.Time `gorm:"type:date;example:'User birthday, format YYYY-MM-DD'" json:"birthday"`
	AvatarURL string     `gorm:"type:varchar(255);example:Avatar URL" json:"avatar_url"`
}

// TableName specifies the MySQL table name to which the struct maps.
func (m *UserM) TableName() string {
	return "user"
}

// BeforeCreate encrypts the plain text password before creating a database record.
func (m *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	m.Password, err = auth.Encrypt(m.Password)
	if err != nil {
		return err
	}
	m.UserID = gid.GenShortID()
	return nil
}
