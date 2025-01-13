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

// OrderLogStore 定义了 comment 模块在 store 层所实现的方法.
type OrderLogStore interface {
	Create(ctx context.Context, orderLog *model.OrderLogM) error
}

// OrderLogStore 接口的实现.
type orderLogs struct {
	db *gorm.DB
}

// 确保 orderLogs 实现了 OrderLogStore 接口.
var _ OrderLogStore = (*orderLogs)(nil)

func newOrderLogs(db *gorm.DB) *orderLogs {
	return &orderLogs{db}
}

// Create 插入一条记录.
func (u *orderLogs) Create(ctx context.Context, orderLog *model.OrderLogM) error {
	return u.db.WithContext(ctx).Create(&orderLog).Error
}
