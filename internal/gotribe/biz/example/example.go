// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package example

import (
	"context"
	"errors"

	"gotribe/internal/gotribe/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ExampleBiz defines functions used to handle example request.
type ExampleBiz interface {
	Create(ctx context.Context, username string, r *v1.CreateExampleRequest) (*v1.CreateExampleResponse, error)
	Update(ctx context.Context, username, exampleID string, r *v1.UpdateExampleRequest) error
	Delete(ctx context.Context, username, exampleID string) error
	DeleteCollection(ctx context.Context, username string, exampleIDs []string) error
	Get(ctx context.Context, username, exampleID string) (*v1.GetExampleResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListExampleResponse, error)
}

// The implementation of ExampleBiz interface.
type exampleBiz struct {
	ds store.IStore
}

// Make sure that exampleBiz implements the ExampleBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ ExampleBiz = (*exampleBiz)(nil)

func New(ds store.IStore) *exampleBiz {
	return &exampleBiz{ds: ds}
}

// CreateTx 实现事务的示例
func (b *exampleBiz) CreateTx(ctx context.Context, username string, r *v1.CreateExampleRequest) error {
	err := b.ds.TX(ctx, func(ctx context.Context) error {
		var exampleM model.ExampleM
		_ = copier.Copy(&exampleM, r)
		exampleM.Username = username

		if err := b.ds.Examples().Create(ctx, &exampleM); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Create is the implementation of the `Create` method in ExampleBiz interface.
func (b *exampleBiz) Create(ctx context.Context, username string, r *v1.CreateExampleRequest) (*v1.CreateExampleResponse, error) {
	var exampleM model.ExampleM
	_ = copier.Copy(&exampleM, r)
	exampleM.Username = username

	if err := b.ds.Examples().Create(ctx, &exampleM); err != nil {
		return nil, err
	}

	return &v1.CreateExampleResponse{ExampleID: exampleM.ExampleID}, nil
}

// Delete is the implementation of the `Delete` method in ExampleBiz interface.
func (b *exampleBiz) Delete(ctx context.Context, username, exampleID string) error {
	if err := b.ds.Examples().Delete(ctx, username, []string{exampleID}); err != nil {
		return err
	}

	return nil
}

// DeleteCollection is the implementation of the `DeleteCollection` method in ExampleBiz interface.
func (b *exampleBiz) DeleteCollection(ctx context.Context, username string, exampleIDs []string) error {
	if err := b.ds.Examples().Delete(ctx, username, exampleIDs); err != nil {
		return err
	}

	return nil
}

// Get is the implementation of the `Get` method in ExampleBiz interface.
func (b *exampleBiz) Get(ctx context.Context, username, exampleID string) (*v1.GetExampleResponse, error) {
	example, err := b.ds.Examples().Get(ctx, username, exampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrExampleNotFound
		}

		return nil, err
	}

	var resp v1.GetExampleResponse
	_ = copier.Copy(&resp, example)

	resp.CreatedAt = example.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = example.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// Update is the implementation of the `Update` method in ExampleBiz interface.
func (b *exampleBiz) Update(ctx context.Context, username, exampleID string, r *v1.UpdateExampleRequest) error {
	exampleM, err := b.ds.Examples().Get(ctx, username, exampleID)
	if err != nil {
		return err
	}

	if r.Title != nil {
		exampleM.Title = *r.Title
	}

	if r.Content != nil {
		exampleM.Content = *r.Content
	}

	if err := b.ds.Examples().Update(ctx, exampleM); err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `List` method in ExampleBiz interface.
func (b *exampleBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListExampleResponse, error) {
	count, list, err := b.ds.Examples().List(ctx, username, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list examples from storage", "err", err)
		return nil, err
	}

	examples := make([]*v1.ExampleInfo, 0, len(list))
	for _, item := range list {
		example := item
		examples = append(examples, &v1.ExampleInfo{
			Username:  example.Username,
			ExampleID: example.ExampleID,
			Title:     example.Title,
			Content:   example.Content,
			CreatedAt: example.CreatedAt.Format(known.TimeFormat),
			UpdatedAt: example.UpdatedAt.Format(known.TimeFormat),
		})
	}

	return &v1.ListExampleResponse{TotalCount: count, Examples: examples}, nil
}
