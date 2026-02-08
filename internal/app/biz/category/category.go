// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package category

import (
	"context"
	"errors"
	v1 "gotribe/pkg/api/v1"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// CategoryBiz defines functions used to handle category request.
type CategoryBiz interface {
	Get(ctx context.Context, categoryID string) (*v1.GetCategoryResponse, error)
	GetChildren(ctx context.Context, categoryID string) (*v1.GetCategoryChildrenResponse, error)
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

// GetChildren is the implementation of the `GetChildren` method in CategoryBiz interface.
func (b *categoryBiz) GetChildren(ctx context.Context, categoryID string) (*v1.GetCategoryChildrenResponse, error) {
	// 先获取到父分类的信息
	category, err := b.ds.Categories().Get(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrCategoryNotFound
		}

		return nil, err
	}
	// 再去获取父类下所有子分类
	categories, err := b.ds.Categories().GetChildren(ctx, category.ID)
	if err != nil {
		return nil, err
	}

	var resp v1.GetCategoryChildrenResponse
	resp.CategoryID = category.CategoryID
	resp.Children = make([]v1.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		var ci v1.CategoryInfo
		_ = copier.Copy(&ci, c)
		ci.CreatedAt = c.CreatedAt.Format(known.TimeFormat)
		ci.UpdatedAt = c.UpdatedAt.Format(known.TimeFormat)
		resp.Children = append(resp.Children, ci)
	}

	// 如果没有子分类，返回一个空切片而不是 null
	if resp.Children == nil {
		resp.Children = []v1.CategoryInfo{}
	}

	return &resp, nil
}
