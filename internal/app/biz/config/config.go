// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package config

import (
	"context"
	"errors"
	"gotribe/pkg/api/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
)

// ConfigBiz defines functions used to handle config request.
type ConfigBiz interface {
	Get(ctx context.Context, alias string) (*v1.GetConfigResponse, error)
}

// The implementation of ConfigBiz interface.
type configBiz struct {
	ds store.IStore
}

// Make sure that configBiz implements the ConfigBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ ConfigBiz = (*configBiz)(nil)

func New(ds store.IStore) *configBiz {
	return &configBiz{ds: ds}
}

// Get is the implementation of the `Get` method in ConfigBiz interface.
func (b *configBiz) Get(ctx context.Context, alias string) (*v1.GetConfigResponse, error) {
	config, err := b.ds.Configs().Get(ctx, alias)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrConfigNotFound
		}

		return nil, err
	}

	var resp v1.GetConfigResponse
	_ = copier.Copy(&resp, config)

	resp.CreatedAt = config.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = config.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}
