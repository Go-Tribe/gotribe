// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gotribe/pkg/api/v1"

	"gorm.io/gorm"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
)

// ColumnStore 定义了 column 模块在 store 层所实现的方法.
type ColumnStore interface {
	Get(ctx context.Context, columnID string) (*model.ColumnM, error)
	List(ctx context.Context, r *v1.ListColumnRequest) (count int64, ret []*model.ColumnM, err error)
}

// ColumnStore 接口的实现.
type columns struct {
	db *gorm.DB
}

// 确保 columns 实现了 ColumnStore 接口.
var _ ColumnStore = (*columns)(nil)

func newColumns(db *gorm.DB) *columns {
	return &columns{db}
}

// Get 根据 columnID 查询指定的 column 数据库记录.
func (u *columns) Get(ctx context.Context, columnID string) (*model.ColumnM, error) {
	var column model.ColumnM
	if err := u.db.WithContext(ctx).Where("column_id = ? and status = ?", columnID, known.STATUS_OK).First(&column).Error; err != nil {
		return nil, err
	}

	return &column, nil
}

func (u *columns) List(ctx context.Context, r *v1.ListColumnRequest) (count int64, ret []*model.ColumnM, err error) {
	// 声明一个空的 []interface{} 切片用于存放查询条件
	queryWhere := make([]interface{}, 0)
	// 逐个创建查询条件并追加到 queryWhere 切片中
	queryWhere = append(queryWhere, []interface{}{"project_id", ctx.Value(known.XProjectIDKey).(string)})
	queryWhere = append(queryWhere, []interface{}{"status", known.STATUS_OK})
	db, err := buildQueryList(u.db, queryWhere, "*", "id desc", r.Offset, r.Limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list column from build where", "err", err)
		return
	}
	// Offset(-1).Limit(-1).Count(&count) 不能少。否则会count=0
	err = db.WithContext(ctx).Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}
