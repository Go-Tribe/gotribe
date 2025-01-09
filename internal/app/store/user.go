// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"errors"
	"gotribe/pkg/api/v1"

	"gorm.io/gorm"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) (*model.UserM, error)
	Get(ctx context.Context, userWhere v1.UserWhere) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	List(ctx context.Context, offset, limit int, userWhere v1.UserWhere) (int64, []*model.UserM, error)
	Delete(ctx context.Context, username string) error
	ListInUserID(ctx context.Context, userID []string) ([]*model.UserM, error)
}

// UserStore 接口的实现.
type users struct {
	db *gorm.DB
}

// 确保 users 实现了 UserStore 接口.
var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// Create 插入一条 user 记录.
func (u *users) Create(ctx context.Context, user *model.UserM) (*model.UserM, error) {
	user.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	result := u.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	// 返回创建好的用户信息
	return user, nil
}

// Get 根据用户名查询指定 user 的数据库记录.
func (u *users) Get(ctx context.Context, userWhere v1.UserWhere) (*model.UserM, error) {
	var user model.UserM
	db, err := buildWhere(u.db, userWhere)
	if err != nil {
		log.C(ctx).Errorw("Failed to Get user from build where", "err", err)
		return nil, err
	}
	if err := db.First(&user).Error; err != nil {
		log.C(ctx).Errorw("Failed to Get user from sql", "err", err)
		return nil, err
	}

	return &user, nil
}

// Update 更新一条 user 数据库记录.
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	log.C(ctx).Infow("user update", "user:", user)
	return u.db.Save(user).Error
}

// List 根据 offset 和 limit 返回 user 列表.
func (u *users) List(ctx context.Context, offset, limit int, userWhere v1.UserWhere) (count int64, ret []*model.UserM, err error) {
	db, err := buildQueryList(u.db, userWhere, "*", "id desc", offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from build where", "err", err)
		return
	}
	// Offset(-1).Limit(-1).Count(&count) 不能少。否则会count=0
	err = db.Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// Delete 根据 username 删除数据库 user 记录.
func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.C(ctx).Errorw("Failed to Delete users from sql", "err", err)
		return err
	}
	return nil
}

func (c *users) ListInUserID(ctx context.Context, userID []string) ([]*model.UserM, error) {
	var usersM []*model.UserM
	if err := c.db.WithContext(ctx).Where("user_id IN ?", userID).Find(&usersM).Error; err != nil {
		return nil, err
	}
	return usersM, nil
}
