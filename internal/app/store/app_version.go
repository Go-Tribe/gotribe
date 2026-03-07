// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"errors"
	"time"

	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// AppVersionStore 定义 app_version 在 store 层的方法.
type AppVersionStore interface {
	// GetLatestRelease 按产品+平台取当前最新版本：release_date <= 当前时间、status=有效，按 release_date 降序取第一条.
	GetLatestRelease(ctx context.Context, productName, platform string) (*model.AppVersionM, error)
}

type appVersions struct {
	db *gorm.DB
}

var _ AppVersionStore = (*appVersions)(nil)

func newAppVersions(db *gorm.DB) *appVersions {
	return &appVersions{db: db}
}

func (s *appVersions) GetLatestRelease(ctx context.Context, productName, platform string) (*model.AppVersionM, error) {
	var m model.AppVersionM
	now := time.Now()
	err := s.db.WithContext(ctx).
		Where("product_name = ? AND platform = ? AND status = ?", productName, platform, known.STATUS_OK).
		Where("release_date IS NOT NULL AND release_date <= ?", now).
		Order("release_date DESC").
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
