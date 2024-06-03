// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"errors"

	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// ExampleStore 定义了 example 模块在 store 层所实现的方法.
type ExampleStore interface {
	Create(ctx context.Context, example *model.ExampleM) error
	Get(ctx context.Context, username, exampleID string) (*model.ExampleM, error)
	Update(ctx context.Context, example *model.ExampleM) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.ExampleM, error)
	Delete(ctx context.Context, username string, exampleIDs []string) error
}

// ExampleStore 接口的实现.
type examples struct {
	db *gorm.DB
}

// 确保 examples 实现了 ExampleStore 接口.
var _ ExampleStore = (*examples)(nil)

func newExamples(db *gorm.DB) *examples {
	return &examples{db}
}

// Create 插入一条 example 记录.
func (u *examples) Create(ctx context.Context, example *model.ExampleM) error {
	return u.db.Create(&example).Error
}

// Get 根据 exampleID 查询指定用户的 example 数据库记录.
func (u *examples) Get(ctx context.Context, username, exampleID string) (*model.ExampleM, error) {
	var example model.ExampleM
	if err := u.db.Where("username = ? and example_id = ?", username, exampleID).First(&example).Error; err != nil {
		return nil, err
	}

	return &example, nil
}

// Update 更新一条 example 数据库记录.
func (u *examples) Update(ctx context.Context, example *model.ExampleM) error {
	return u.db.Save(example).Error
}

// List 根据 offset 和 limit 返回指定用户的 example 列表.
func (u *examples) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.ExampleM, err error) {
	err = u.db.Where("username = ?", username).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete 根据 username, exampleID 删除数据库 example 记录.
func (u *examples) Delete(ctx context.Context, username string, exampleIDs []string) error {
	err := u.db.Where("username = ? and example_id in (?)", username, exampleIDs).Delete(&model.ExampleM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
