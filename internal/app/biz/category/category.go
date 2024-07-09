// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package category

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

// CategoryBiz defines functions used to handle category request.
type CategoryBiz interface {
	Get(ctx context.Context, categoryID string) (*v1.GetCategoryResponse, error)
}

// The implementation of CategoryBiz interface.
type categoryBiz struct {
	ds store.IStore
}

// Make sure that categoryBiz implements the CategoryBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ CategoryBiz = (*categoryBiz)(nil)

func New(ds store.IStore) *categoryBiz {
	return &categoryBiz{ds: ds}
}

// Get is the implementation of the `Get` method in CategoryBiz interface.
func (b *categoryBiz) Get(ctx context.Context, categoryID string) (*v1.GetCategoryResponse, error) {
	category, err := b.ds.Categories().Get(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrCategoryNotFound
		}

		return nil, err
	}

	var resp v1.GetCategoryResponse
	_ = copier.Copy(&resp, category)

	resp.CreatedAt = category.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = category.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}
