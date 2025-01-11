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

// PointAvailableStore 定义了 comment 模块在 store 层所实现的方法.
type PointAvailableStore interface {
	SumPoints(ctx context.Context, userID, projectID string) (float64, error)
	Create(ctx context.Context, pointAvailable *model.PointAvailableM) error
}

// PointAvailableStore 接口的实现.
type pointAvailables struct {
	db *gorm.DB
}

// 确保 pointAvailable 实现了 PointAvailableStore 接口.
var _ PointAvailableStore = (*pointAvailables)(nil)

func newPointAvailables(db *gorm.DB) *pointAvailables {
	return &pointAvailables{db}
}

// SumPoints 根据 userID 求和指定用户的 points 数据库记录.
func (u *pointAvailables) SumPoints(ctx context.Context, userID, projectID string) (float64, error) {
	var sumPoints float64
	result := u.db.WithContext(ctx).Model(&model.PointAvailableM{}).Select("sum(points)").
		Where("user_id = ? and project_id = ?", userID, projectID).
		Scan(&sumPoints)
	if result.Error != nil {
		return 0, result.Error
	}
	return sumPoints, nil
}

func (u *pointAvailables) Create(ctx context.Context, pointAvailable *model.PointAvailableM) error {
	return u.db.WithContext(ctx).Create(&pointAvailable).Error
}
