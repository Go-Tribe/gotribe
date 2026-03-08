// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package app_version

import (
	"context"
	"errors"
	"strconv"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"
	"gotribe/pkg/api/v1"

	"gorm.io/gorm"
)

// AppVersionBiz 定义 app_version 在 biz 层的方法.
type AppVersionBiz interface {
	// GetByAppVersionID 根据 app_version_id 查版本记录，用于从表里取 clientVersionCode.
	GetByAppVersionID(ctx context.Context, appVersionID string) (*model.AppVersionM, error)
	// GetLatestRelease 根据 ClientName+OS+OSArch 获取最新发布版本，并根据传入的 clientVersionCode 判断是否需强制升级.
	GetLatestRelease(ctx context.Context, clientName, os, osArch string, clientVersionCode int) (*v1.GetLatestReleaseResponse, error)
}

type appVersionBiz struct {
	ds store.IStore
}

var _ AppVersionBiz = (*appVersionBiz)(nil)

func New(ds store.IStore) *appVersionBiz {
	return &appVersionBiz{ds: ds}
}

func (b *appVersionBiz) GetByAppVersionID(ctx context.Context, appVersionID string) (*model.AppVersionM, error) {
	return b.ds.AppVersions().GetByAppVersionID(ctx, appVersionID)
}

func (b *appVersionBiz) GetLatestRelease(ctx context.Context, clientName, os, osArch string, clientVersionCode int) (*v1.GetLatestReleaseResponse, error) {
	m, err := b.ds.AppVersions().GetLatestRelease(ctx, clientName, os, osArch)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPageNotFound
		}
		return nil, err
	} else if m == nil {
		return nil, errno.ErrInvalidParameter
	}

	resp := &v1.GetLatestReleaseResponse{
		ProductName:             m.ClientName,
		Platform:                m.OS,
		VersionCode:             m.ClientVersionCode,
		VersionName:             m.ClientVersion,
		MinSupportedVersionCode: m.MinSupportedVersionCode,
		ForceUpdate:             m.ForceUpdate,
		Title:                   m.Title,
		Content:                 m.Content,
		DownloadURL:             m.DownloadURL,
		FileSize:                m.FileSize,
	}
	if m.ReleaseDate != nil {
		resp.ReleaseDate = m.ReleaseDate.Format(known.TimeFormat)
	}

	// 是否需要强制升级：最新版标记了强制 或 客户端版本号 < 最新版要求的最低兼容版本号
	resp.NeedForceUpdate = m.ForceUpdate == 1 || (clientVersionCode < m.MinSupportedVersionCode)

	return resp, nil
}

// ParseClientVersionCode 从字符串解析客户端版本号，解析失败返回 0.
func ParseClientVersionCode(s string) int {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
