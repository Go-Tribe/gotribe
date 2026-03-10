// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"time"

	"gorm.io/gorm"
)

// VerificationCodeM 验证码记录（邮件/短信等）.
type VerificationCodeM struct {
	gorm.Model
	Channel   string    `gorm:"type:varchar(20);not null;comment:渠道类型：email" json:"channel"`
	Trigger   string    `gorm:"type:varchar(20);not null;comment:场景：register" json:"trigger"`
	Target    string    `gorm:"type:varchar(128);not null;index:idx_target_trigger;comment:接收方，如邮箱" json:"target"`
	Code      string    `gorm:"type:varchar(20);not null;comment:验证码" json:"code"`
	UserID    string    `gorm:"type:varchar(64);index;comment:用户ID" json:"userID"`
	ProjectID string    `gorm:"type:varchar(64);comment:项目ID" json:"projectID"`
	ExpireAt  time.Time `gorm:"not null;comment:过期时间" json:"expireAt"`
	Verified  int       `gorm:"type:smallint;default:0;comment:是否已验证：0-未验证，1-已验证" json:"verified"`
	TraceID   string    `gorm:"type:varchar(64);comment:追踪ID，用于关联发送链路" json:"traceID"`
}

// TableName 指定表名
func (VerificationCodeM) TableName() string {
	return "verification_code"
}
