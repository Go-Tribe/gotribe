// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// ProductStore 定义了 product 模块在 store 层所实现的方法.
type ProductStore interface {
	Get(ctx context.Context, productID string) (*model.ProductM, error)
	List(ctx context.Context, categoryID string, offset, limit int) (int64, []*model.ProductM, error)
}

// ProductStore 接口的实现.
type products struct {
	db *gorm.DB
}

// 确保 products 实现了 ProductStore 接口.
var _ ProductStore = (*products)(nil)

func newProducts(db *gorm.DB) *products {
	return &products{db}
}

// Get 根据 productID 查询指定用户的 product 数据库记录.
func (u *products) Get(ctx context.Context, productID string) (*model.ProductM, error) {
	var product model.ProductM
	if err := u.db.WithContext(ctx).Where("product_id = ? and enable = ?", productID, known.STATUS_PUBLIC).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// List 根据 offset 和 limit 返回指定分类的 product 列表.
func (u *products) List(ctx context.Context, categoryID string, offset, limit int) (count int64, ret []*model.ProductM, err error) {
	err = u.db.WithContext(ctx).Where("category_id = ? and enable = ?", categoryID, known.STATUS_PUBLIC).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}
