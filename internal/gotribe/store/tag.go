// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"

	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// TagStore 定义了 example 模块在 store 层所实现的方法.
type TagStore interface {
	Get(ctx context.Context, tagID string) (*model.TagM, error)
	GetTags(ctx context.Context, tagIDs []string) (ret []*model.TagM, err error)
}

// TagStore 接口的实现.
type tags struct {
	db *gorm.DB
}

// 确保 tag 实现了 TagStore 接口.
var _ TagStore = (*tags)(nil)

func newTags(db *gorm.DB) *tags {
	return &tags{db}
}

// Get 根据 exampleID 查询指定用户的 example 数据库记录.
func (u *tags) Get(ctx context.Context, tagID string) (*model.TagM, error) {
	var tag model.TagM
	if err := u.db.Where("tag_id = ?", tagID).First(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// GetTags.
func (u *tags) GetTags(ctx context.Context, tagIDs []string) (ret []*model.TagM, err error) {
	err = u.db.Where("tag_id in (?)", tagIDs).Find(&ret).Error
	return
}
