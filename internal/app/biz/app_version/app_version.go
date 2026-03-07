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
	"gotribe/pkg/api/v1"

	"gorm.io/gorm"
)

// AppVersionBiz 定义 app_version 在 biz 层的方法.
type AppVersionBiz interface {
	// GetLatestRelease 获取当前产品+平台的最新发布版本，并计算 needForceUpdate.
	// clientVersionCode 来自请求头 x-platform-version-code，解析失败时按 0 处理（倾向于需要强制升级）.
	GetLatestRelease(ctx context.Context, productName, platform string, clientVersionCode int) (*v1.GetLatestReleaseResponse, error)
}

type appVersionBiz struct {
	ds store.IStore
}

var _ AppVersionBiz = (*appVersionBiz)(nil)

func New(ds store.IStore) *appVersionBiz {
	return &appVersionBiz{ds: ds}
}

func (b *appVersionBiz) GetLatestRelease(ctx context.Context, productName, platform string, clientVersionCode int) (*v1.GetLatestReleaseResponse, error) {
	m, err := b.ds.AppVersions().GetLatestRelease(ctx, productName, platform)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPageNotFound
		}
		return nil, err
	} else if m == nil {
		return nil, errno.ErrInvalidParameter
	}

	resp := &v1.GetLatestReleaseResponse{
		ProductName:             m.ProductName,
		Platform:                m.Platform,
		VersionCode:             m.VersionCode,
		VersionName:             m.VersionName,
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
