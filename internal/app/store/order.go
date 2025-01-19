// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gorm.io/gorm"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/v1"
)

// OrderStore 定义了 comment 模块在 store 层所实现的方法.
type OrderStore interface {
	Create(ctx context.Context, order *model.OrderM) (*model.OrderM, error)
	Get(ctx context.Context, orderWhere v1.OrderWhere) (*model.OrderM, error)
	List(ctx context.Context, offset, limit int, orderWhere v1.OrderWhere) (count int64, ret []*model.OrderM, err error)
	Update(ctx context.Context, order *model.OrderM) error
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

func (u *orders) Get(ctx context.Context, orderWhere v1.OrderWhere) (*model.OrderM, error) {
	var order model.OrderM
	db, err := buildWhere(u.db, orderWhere)
	if err != nil {
		log.C(ctx).Errorw("Failed to Get user from build where", "err", err)
		return nil, err
	}
	if err := db.WithContext(ctx).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (u *orders) List(ctx context.Context, offset, limit int, orderWhere v1.OrderWhere) (count int64, ret []*model.OrderM, err error) {
	db, err := buildQueryList(u.db, orderWhere, "*", "id desc", offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list order from build where", "err", err)
		return
	}
	err = db.WithContext(ctx).Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}

func (u *orders) Update(ctx context.Context, order *model.OrderM) error {
	log.C(ctx).Infow("order update", "order:", order)
	return u.db.WithContext(ctx).Save(order).Error
}
