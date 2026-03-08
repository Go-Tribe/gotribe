// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"time"

	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ConversationLogM struct {
	gorm.Model
	SID              string     `gorm:"type:char(10);not null;default:'';comment:短id" json:"sid"`
	UserID           string     `gorm:"type:char(10);not null;default:'';Index;comment:用户ID" json:"userID"`
	Username         string     `gorm:"type:varchar(255);not null;default:'';Index;comment:用户名" json:"username"`
	Question         string     `gorm:"type:varchar(500);not null;default:'';comment:用户问题" json:"question"`
	Answer           string     `gorm:"type:varchar(500);not null;default:'';comment:AI回答" json:"answer"`
	ModelName        string     `gorm:"type:varchar(100);not null;default:'';comment:模型名称" json:"modelName"`
	PromptTokens     int        `gorm:"not null;default:0;comment:输入token数" json:"promptTokens"`
	CacheTokens      int        `gorm:"not null;default:0;comment:缓存输入token数" json:"cacheTokens"`
	CompletionTokens int        `gorm:"not null;default:0;comment:输出token数" json:"completionTokens"`
	TotalTokens      int        `gorm:"not null;default:0;comment:总token数" json:"totalTokens"`
	PointsCost       float64    `gorm:"type:decimal(10,2);not null;default:0;comment:消耗点数(精确到小数点2位)" json:"pointsCost"`
	StartTime        time.Time  `gorm:"type:timestamp;not null;comment:开始时间(毫秒精度)" json:"startTime"`
	EndTime          *time.Time `gorm:"type:timestamp;comment:结束时间" json:"endTime,omitempty"`
	TotalTimeMs      int        `gorm:"not null;default:0;comment:总耗时(毫秒)" json:"totalTimeMs"`
	DeductionStatus  int        `gorm:"type:smallint;not null;default:1;comment:扣费状态(1-待扣费,2-扣费成功,3-扣费失败,4-无需扣费)" json:"deductionStatus"`
	DeductionType    int        `gorm:"type:smallint;not null;default:1;comment:扣费类型(1-点数扣费,2-内部白名单用户)" json:"deductionType"`
	DeductionTime    *time.Time `gorm:"type:timestamp;comment:扣费时间" json:"deductionTime,omitempty"`
	AppVersionId     string     `gorm:"type:varchar(50);not null;default:'';comment:应用版本id" json:"appVersionId"`
	ClientName       string     `gorm:"type:varchar(50);not null;default:'';comment:客户端名" json:"clientName"`
	ClientVersion    string     `gorm:"type:varchar(30);not null;default:'';comment:客户端版本" json:"clientVersion"`
	DeviceModel      string     `gorm:"type:varchar(30);not null;default:'';comment:设备号" json:"deviceModel"`
	OS               string     `gorm:"type:varchar(30);not null;default:'';comment:ios/android/darwin/windows/linux" json:"os"`
	OSVersion        string     `gorm:"type:varchar(30);not null;default:'';comment:系统版本" json:"osVersion"`
	OSArch           string     `gorm:"type:varchar(30);not null;default:'';comment:系统架构，如：arm64" json:"osArch"`
	CPUModel         string     `gorm:"type:varchar(30);not null;default:'';comment:如：Intel Core i9-14900HX" json:"cpuModel"`
	RequestID        string     `gorm:"type:varchar(64);not null;default:'';comment:请求ID(用于幂等)" json:"requestId"`
	IPAddress        string     `gorm:"type:varchar(64);not null;default:'';comment:客户端IP" json:"ipAddress"`
	Status           int        `gorm:"type:smallint;not null;default:1;comment:状态(1-成功,2-失败)" json:"status"`
	ErrorMessage     string     `gorm:"type:varchar(500);default:'';comment:错误信息" json:"errorMessage"`
	Ext              string     `gorm:"type:text;comment:扩展信息(JSON格式)" json:"ext"`
}

func (m *ConversationLogM) TableName() string {
	return "conversation_log"
}

func (m *ConversationLogM) BeforeCreate(tx *gorm.DB) error {
	m.SID = gid.GenShortID()
	return nil
}
