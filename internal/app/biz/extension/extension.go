// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package extension

import (
	"context"
	"errors"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/pkg/api/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ExtensionBiz 定义了 extension 模块在 biz 层所实现的方法.
type ExtensionBiz interface {
	Get(ctx context.Context, extensionID string) (*v1.GetExtensionResponse, error)
	MarketList(ctx context.Context, r *v1.ListExtensionRequest) (*v1.ListExtensionResponse, error)
}

type extensionBiz struct {
	ds store.IStore
}

var _ ExtensionBiz = (*extensionBiz)(nil)

// New 创建一个实现了 ExtensionBiz 接口的实例.
func New(ds store.IStore) *extensionBiz {
	return &extensionBiz{ds: ds}
}

// Get 根据 extensionID 查询插件详情.
func (b *extensionBiz) Get(ctx context.Context, extensionID string) (*v1.GetExtensionResponse, error) {
	m, err := b.ds.Extensions().Get(ctx, extensionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPageNotFound
		}
		return nil, err
	}
	var resp v1.GetExtensionResponse
	_ = copier.Copy(&resp, m)
	resp.CreatedAt = m.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = m.UpdatedAt.Format(known.TimeFormat)
	return &resp, nil
}

// MarketList 返回插件列表，可选按 categoryID、type 筛选.
func (b *extensionBiz) MarketList(ctx context.Context, r *v1.ListExtensionRequest) (*v1.ListExtensionResponse, error) {
	count, list, err := b.ds.Extensions().List(ctx, r.CategoryID, r.Type, r.Offset, r.Limit)
	if err != nil {
		return nil, err
	}
	extensions := make([]*v1.ExtensionInfo, 0, len(list))
	for _, m := range list {
		var info v1.ExtensionInfo
		_ = copier.Copy(&info, m)
		info.CreatedAt = m.CreatedAt.Format(known.TimeFormat)
		info.UpdatedAt = m.UpdatedAt.Format(known.TimeFormat)
		extensions = append(extensions, &info)
	}
	return &v1.ListExtensionResponse{TotalCount: count, Extensions: extensions}, nil
}
