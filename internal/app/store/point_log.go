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

// PointLogStore 定义了 comment 模块在 store 层所实现的方法.
type PointLogStore interface {
	SumPoints(ctx context.Context, userID, projectID string) (float64, error)
	Create(ctx context.Context, pointLog *model.PointLogM) (uint, error)
}

// PointLogStore 接口的实现.
type pointLogs struct {
	db *gorm.DB
}

// 确保 pointLog 实现了 PointLogStore 接口.
var _ PointLogStore = (*pointLogs)(nil)

func newPointLogs(db *gorm.DB) *pointLogs {
	return &pointLogs{db}
}

// SumPoints 根据 userID 求和指定用户的 points 数据库记录.
func (u *pointLogs) SumPoints(ctx context.Context, userID, projectID string) (float64, error) {
	var sumPoints float64
	result := u.db.WithContext(ctx).Model(&model.PointLogM{}).Select("sum(points)").
		Where("user_id = ? and project_id = ?", userID, projectID).
		Scan(&sumPoints)
	if result.Error != nil {
		return 0, result.Error
	}
	return sumPoints, nil
}

func (u *pointLogs) Create(ctx context.Context, pointLog *model.PointLogM) (uint, error) {
	result := u.db.WithContext(ctx).Create(pointLog)
	if result.Error != nil {
		return 0, result.Error
	}
	return pointLog.ID, nil
}
