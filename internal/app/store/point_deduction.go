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

// PointDeductionStore 定义了 comment 模块在 store 层所实现的方法.
type PointDeductionStore interface {
	Create(ctx context.Context, pointDeduction *model.PointDeductionM) error
}

// PointDeductionStore 接口的实现.
type pointDeductions struct {
	db *gorm.DB
}

// 确保 pointDeduction 实现了 PointDeductionStore 接口.
var _ PointDeductionStore = (*pointDeductions)(nil)

func newPointDeductions(db *gorm.DB) *pointDeductions {
	return &pointDeductions{db}
}

func (u *pointDeductions) Create(ctx context.Context, pointDeduction *model.PointDeductionM) error {
	return u.db.WithContext(ctx).Create(&pointDeduction).Error
}
