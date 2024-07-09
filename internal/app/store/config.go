// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"

	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// ConfigStore 定义了 config 模块在 store 层所实现的方法.
type ConfigStore interface {
	Get(ctx context.Context, alias string) (*model.ConfigM, error)
}

// ConfigStore 接口的实现.
type configs struct {
	db *gorm.DB
}

// 确保 configs 实现了 ConfigStore 接口.
var _ ConfigStore = (*configs)(nil)

func newConfigs(db *gorm.DB) *configs {
	return &configs{db}
}

// Get 根据 configID 查询指定用户的 config 数据库记录.
func (u *configs) Get(ctx context.Context, alias string) (*model.ConfigM, error) {
	var config model.ConfigM
	if err := u.db.Where("alias = ?", alias).First(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}
