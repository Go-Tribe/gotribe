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

// ProjectStore 定义了 example 模块在 store 层所实现的方法.
type ProjectStore interface {
	Get(ctx context.Context, name string) (*model.ProjectM, error)
}

// ProjectStore 接口的实现.
type projects struct {
	db *gorm.DB
}

// 确保 project 实现了 ProjectStore 接口.
var _ ProjectStore = (*projects)(nil)

func newProjects(db *gorm.DB) *projects {
	return &projects{db}
}

// Get 根据 exampleID 查询指定用户的 example 数据库记录.
func (u *projects) Get(ctx context.Context, name string) (*model.ProjectM, error) {
	var project model.ProjectM
	if err := u.db.Where("name = ?", name).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}
