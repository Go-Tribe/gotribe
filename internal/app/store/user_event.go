// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gorm.io/gorm"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"
)

// UserEventStore 定义了 userEvent 模块在 store 层所实现的方法.
type UserEventStore interface {
	Create(ctx context.Context, userEvent *model.UserEventM) (*model.UserEventM, error)
}

// UserEventStore 接口的实现.
type userEvents struct {
	db *gorm.DB
}

// 确保 userEvents 实现了 UserEventStore 接口.
var _ UserEventStore = (*userEvents)(nil)

func newUserEvents(db *gorm.DB) *userEvents {
	return &userEvents{db}
}

// Create 插入一条 userEvent 记录.
func (u *userEvents) Create(ctx context.Context, userEvent *model.UserEventM) (*model.UserEventM, error) {
	userEvent.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	result := u.db.WithContext(ctx).Create(&userEvent)
	if result.Error != nil {
		return nil, result.Error
	}
	// 返回创建好的用户信息
	return userEvent, nil
}
