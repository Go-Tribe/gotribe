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

// FeedBackStore 定义了 comment 模块在 store 层所实现的方法.
type FeedBackStore interface {
	Create(ctx context.Context, feedBack *model.FeedbackM) error
}

// FeedBackStore 接口的实现.
type feedBacks struct {
	db *gorm.DB
}

// 确保 feedBacks 实现了 FeedBackStore 接口.
var _ FeedBackStore = (*feedBacks)(nil)

func newFeedBacks(db *gorm.DB) *feedBacks {
	return &feedBacks{db}
}

// Create 插入一条 comment 记录.
func (u *feedBacks) Create(ctx context.Context, feedBack *model.FeedbackM) error {
	return u.db.WithContext(ctx).Create(&feedBack).Error
}
