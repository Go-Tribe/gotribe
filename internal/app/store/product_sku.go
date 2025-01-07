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

// ProductSKUStore 定义了 productSKU 模块在 store 层所实现的方法.
type ProductSKUStore interface {
	Get(ctx context.Context, productSKUID string) (*model.ProductSKUM, error)
	List(ctx context.Context, productID string) ([]*model.ProductSKUM, error)
}

// ProductSKUStore 接口的实现.
type productSKUs struct {
	db *gorm.DB
}

// 确保 productSKUs 实现了 ProductSKUStore 接口.
var _ ProductSKUStore = (*productSKUs)(nil)

func newProductSKUs(db *gorm.DB) *productSKUs {
	return &productSKUs{db}
}

// Get 根据 productSKUID 查询指定用户的 productSKU 数据库记录.
func (u *productSKUs) Get(ctx context.Context, productSKUID string) (*model.ProductSKUM, error) {
	var productSKU model.ProductSKUM
	if err := u.db.Where("sku_id = ?", productSKUID).First(&productSKU).Error; err != nil {
		return nil, err
	}

	return &productSKU, nil
}

// List 根据 offset 和 limit 返回指定分类的 productSKU 列表.
func (u *productSKUs) List(ctx context.Context, productID string) (ret []*model.ProductSKUM, err error) {
	err = u.db.Where("product_id = ?", productID).Order("id desc").Find(&ret).
		Error

	return
}
