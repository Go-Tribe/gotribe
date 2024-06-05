// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package column

import (
	"context"
	"errors"

	"gotribe/internal/gotribe/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ColumnBiz defines functions used to handle column request.
type ColumnBiz interface {
	Get(ctx context.Context, columnID string) (*v1.GetColumnResponse, error)
	List(ctx context.Context, r *v1.ListColumnRequest) (*v1.LisColumnResponse, error)
}

// The implementation of ColumnBiz interface.
type columnBiz struct {
	ds store.IStore
}

// Make sure that columnBiz implements the ColumnBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ ColumnBiz = (*columnBiz)(nil)

func New(ds store.IStore) *columnBiz {
	return &columnBiz{ds: ds}
}

// Get is the implementation of the `Get` method in ColumnBiz interface.
func (b *columnBiz) Get(ctx context.Context, columnID string) (*v1.GetColumnResponse, error) {
	column, err := b.ds.Columns().Get(ctx, columnID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrColumnNotFound
		}

		return nil, err
	}

	var resp v1.GetColumnResponse
	_ = copier.Copy(&resp, column)

	resp.CreatedAt = column.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = column.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// List is the implementation of the `List` method in colnmnBiz interface.
func (b *columnBiz) List(ctx context.Context, r *v1.ListColumnRequest) (*v1.LisColumnResponse, error) {
	count, list, err := b.ds.Columns().List(ctx, r)
	if err != nil {
		log.C(ctx).Errorw("Failed to list colnmns from storage", "err", err)
		return nil, err
	}

	columns := make([]*v1.ColumnInfo, 0, len(list))
	for _, item := range list {
		column := item
		// 分类信息
		count, postM, err := b.ds.Posts().List(ctx, &v1.ListPostRequest{ColumnID: column.ColumnID, Limit: r.PostLimit, Type: gconvert.String(known.POST_TYPE_POST)})
		if err != nil {
			log.C(ctx).Errorw("Failed to get postM from storage", "err", err)
			return nil, err
		}
		var posts []*v1.PostInfo
		_ = copier.Copy(&posts, postM)

		columns = append(columns, &v1.ColumnInfo{
			Title:       column.Title,
			Description: column.Description,
			Posts:       posts,
			PostCount:   count,
			Icon:        column.Icon,
			ColumnID:    column.ColumnID,
			CreatedAt:   column.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:   column.UpdatedAt.Format(known.TimeFormat),
		})
	}

	return &v1.LisColumnResponse{TotalCount: count, Columns: columns}, nil
}
