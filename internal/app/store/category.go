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

// CategoryStore 定义了 comment 模块在 store 层所实现的方法.
type CategoryStore interface {
	Get(ctx context.Context, categoryID string) (*model.CategoryM, error)
}

// CategoryStore 接口的实现.
type categories struct {
	db *gorm.DB
}

// 确保 category 实现了 CategoryStore 接口.
var _ CategoryStore = (*categories)(nil)

func newCategories(db *gorm.DB) *categories {
	return &categories{db}
}

// Get 根据 exampleID 查询指定用户的 comment 数据库记录.
func (u *categories) Get(ctx context.Context, categoryID string) (*model.CategoryM, error) {
	var category model.CategoryM
	if err := u.db.Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
