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
	// GetByAppVersionID 根据 app_version_id 查一条版本记录.
	GetByAppVersionID(ctx context.Context, appVersionID string) (*model.AppVersionM, error)
	// GetLatestRelease 按 ClientName(产品)+OS+OSArch 取当前最新版本：release_date <= 当前时间、status=有效，按 release_date 降序取第一条.
	GetLatestRelease(ctx context.Context, clientName, os, osArch string) (*model.AppVersionM, error)
}

type appVersions struct {
	db *gorm.DB
}

var _ AppVersionStore = (*appVersions)(nil)

func newAppVersions(db *gorm.DB) *appVersions {
	return &appVersions{db: db}
}

func (s *appVersions) GetByAppVersionID(ctx context.Context, appVersionID string) (*model.AppVersionM, error) {
	if appVersionID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var m model.AppVersionM
	err := s.db.WithContext(ctx).Where("app_version_id = ?", appVersionID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *appVersions) GetLatestRelease(ctx context.Context, clientName, os, osArch string) (*model.AppVersionM, error) {
	var m model.AppVersionM
	now := time.Now()
	db := s.db.WithContext(ctx).
		Where("client_name = ? AND os = ? AND status = ?", clientName, os, known.STATUS_OK).
		Where("release_date IS NOT NULL AND release_date <= ?", now)
	if osArch != "" {
		db = db.Where("(os_arch = ? OR os_arch = '' OR os_arch IS NULL)", osArch)
	}
	err := db.Order("release_date DESC").First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
