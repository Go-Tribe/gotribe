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

// ExtensionStore 定义了 extension 模块在 store 层所实现的方法.
type ExtensionStore interface {
	Get(ctx context.Context, extensionID string) (*model.ExtensionM, error)
	List(ctx context.Context, categoryID string, typeFilter uint, offset, limit int) (int64, []*model.ExtensionM, error)
}

type extensions struct {
	db *gorm.DB
}

var _ ExtensionStore = (*extensions)(nil)

func newExtensions(db *gorm.DB) *extensions {
	return &extensions{db: db}
}

// Get 根据 extensionID 查询插件详情.
func (e *extensions) Get(ctx context.Context, extensionID string) (*model.ExtensionM, error) {
	var m model.ExtensionM
	if err := e.db.WithContext(ctx).Where("extension_id = ?", extensionID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// List 返回插件列表，categoryID 非空或 typeFilter 非 0 时参与筛选.
func (e *extensions) List(ctx context.Context, categoryID string, typeFilter uint, offset, limit int) (int64, []*model.ExtensionM, error) {
	db := e.db.WithContext(ctx).Model(&model.ExtensionM{})
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}
	if typeFilter != 0 {
		db = db.Where("type = ?", typeFilter)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	var list []*model.ExtensionM
	err := db.Order("id desc").
		Limit(defaultLimit(limit)).
		Offset(offset).
		Find(&list).Error
	return count, list, err
}
