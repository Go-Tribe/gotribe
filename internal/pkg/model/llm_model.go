// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type LLMModelM struct {
	gorm.Model
	Sid          string  `gorm:"type:varchar(10);uniqueIndex;not null;comment:短ID(唯一标识)" json:"sid"`
	Name         string  `gorm:"type:varchar(100);not null;comment:名称" json:"name"`
	Description  string  `gorm:"type:varchar(200);default:'';comment:描述" json:"description,omitempty"`
	Icon         string  `gorm:"type:varchar(200);default:'';comment:图标URL或Base64" json:"icon,omitempty"`
	Type         string  `gorm:"type:varchar(20);not null;comment:调用类型(openai-openai风格)" json:"type"`
	ApiKey       string  `gorm:"type:varchar(200);not null;comment:API密钥" json:"apiKey"`
	BaseURL      string  `gorm:"type:varchar(200);not null;comment:API基础地址" json:"baseUrl"`
	ModelName    string  `gorm:"type:varchar(100);not null;comment:模型名" json:"modelName"`
	Params       string  `gorm:"type:varchar(500);default:'';comment:默认参数配置(JSON格式)" json:"params"`
	Header       string  `gorm:"type:varchar(500);default:'';comment:HTTP头配置(JSON格式)" json:"header"`
	ClientHeader string  `gorm:"type:varchar(500);default:'';comment:取客户端头传入的头" json:"clientHeader"`
	Status       int     `gorm:"type:smallint;not null;default:1;comment:状态(1-启用,2-禁用)" json:"status"`
	RateLimit    int     `gorm:"not null;default:0;comment:每分钟请求限制(0表示不限制)" json:"rateLimit"`
	TokenLimit   int     `gorm:"not null;default:0;comment:每分钟Token限制(0表示不限制)" json:"tokenLimit"`
	InputPrice   float64 `gorm:"type:decimal(10,5);not null;default:0;comment:输入每千token消耗点数" json:"inputPrice"`
	CachePrice   float64 `gorm:"type:decimal(10,5);not null;default:0;comment:缓存输入每千token消耗点数" json:"cachePrice"`
	OutputPrice  float64 `gorm:"type:decimal(10,5);not null;default:0;comment:输出每千token消耗点数" json:"outputPrice"`
}

func (m *LLMModelM) TableName() string {
	return "llm_model"
}

func (m *LLMModelM) BeforeCreate(tx *gorm.DB) error {
	m.Sid = gid.GenShortID()
	return nil
}
