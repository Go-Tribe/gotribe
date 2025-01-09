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

// AdStore 定义了 ad 模块在 store 层所实现的方法.
type AdStore interface {
	List(ctx context.Context, sceneID string, offset, limit int) (int64, []*model.AdM, error)
}

// AdStore 接口的实现.
type ads struct {
	db *gorm.DB
}

// 确保 ads 实现了 AdStore 接口.
var _ AdStore = (*ads)(nil)

func newAds(db *gorm.DB) *ads {
	return &ads{db}
}

// List 根据 offset 和 limit 返回指定用户的 ad 列表.
func (u *ads) List(ctx context.Context, sceneID string, offset, limit int) (count int64, ret []*model.AdM, err error) {
	err = u.db.WithContext(ctx).Where("scene_id = ? and status = ?", sceneID, known.STATUS_PUBLIC).Offset(offset).Limit(defaultLimit(limit)).Order("sort desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}
