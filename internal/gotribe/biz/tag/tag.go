// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package tag

import (
	"context"
	"errors"

	"gotribe/internal/gotribe/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// TagBiz defines functions used to handle tag request.
type TagBiz interface {
	Get(ctx context.Context, tagID string) (*v1.GetTagResponse, error)
}

// The implementation of TagBiz interface.
type tagBiz struct {
	ds store.IStore
}

// Make sure that tagBiz implements the TagBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ TagBiz = (*tagBiz)(nil)

func New(ds store.IStore) *tagBiz {
	return &tagBiz{ds: ds}
}

// Get is the implementation of the `Get` method in TagBiz interface.
func (b *tagBiz) Get(ctx context.Context, tagID string) (*v1.GetTagResponse, error) {
	tag, err := b.ds.Tags().Get(ctx, tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrTagNotFound
		}

		return nil, err
	}

	var resp v1.GetTagResponse
	_ = copier.Copy(&resp, tag)

	resp.CreatedAt = tag.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = tag.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}
