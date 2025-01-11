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

// OrderStore 定义了 comment 模块在 store 层所实现的方法.
type OrderStore interface {
	Create(ctx context.Context, order *model.OrderM) (*model.OrderM, error)
}

// OrderStore 接口的实现.
type orders struct {
	db *gorm.DB
}

// 确保 orders 实现了 OrderStore 接口.
var _ OrderStore = (*orders)(nil)

func newOrders(db *gorm.DB) *orders {
	return &orders{db}
}

// Create 插入一条记录.
func (u *orders) Create(ctx context.Context, order *model.OrderM) (*model.OrderM, error) {
	result := u.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return nil, result.Error
	}

	// 假设 order.ID 是自增主键
	var createdOrder model.OrderM
	if err := u.db.WithContext(ctx).First(&createdOrder, order.ID).Error; err != nil {
		return nil, err
	}

	return &createdOrder, nil
}
