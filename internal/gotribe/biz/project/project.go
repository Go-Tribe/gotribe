// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package project

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

// ProjectBiz defines functions used to handle project request.
type ProjectBiz interface {
	Get(ctx context.Context, alias string) (*v1.GetProjectResponse, error)
}

// The implementation of ProjectBiz interface.
type projectBiz struct {
	ds store.IStore
}

// Make sure that projectBiz implements the ProjectBiz interface.
// We can find this problem in the compile sprojecte with the following assignment statement.
var _ ProjectBiz = (*projectBiz)(nil)

func New(ds store.IStore) *projectBiz {
	return &projectBiz{ds: ds}
}

// Get is the implementation of the `Get` method in ProjectBiz interface.
func (b *projectBiz) Get(ctx context.Context, alias string) (*v1.GetProjectResponse, error) {
	project, err := b.ds.Projects().Get(ctx, alias)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrProjectNotFound
		}

		return nil, err
	}

	var resp v1.GetProjectResponse
	_ = copier.Copy(&resp, project)

	resp.CreatedAt = project.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = project.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}
