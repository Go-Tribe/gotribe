// Copyright 2024 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"errors"
	"strings"

	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/gotribe/v1"

	"github.com/dengmengmian/ghelper/gconvert"
	"gorm.io/gorm"
)

// PostStore 定义了 post 模块在 store 层所实现的方法.
type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, query v1.PostQueryParams) (*model.PostM, error)
	Update(ctx context.Context, post *model.PostM) error
	List(ctx context.Context, r *v1.ListPostRequest) (int64, []*model.PostM, error)
	Delete(ctx context.Context, username string, post_ids []string) error
}

// PostStore 接口的实现.
type posts struct {
	db *gorm.DB
}

// 确保 posts 实现了 PostStore 接口.
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

// Create 插入一条 post 记录.
func (u *posts) Create(ctx context.Context, post *model.PostM) error {
	return u.db.Create(&post).Error
}

// Get 根据 post_id 查询指定用户的 post 数据库记录.
func (u *posts) Get(ctx context.Context, query v1.PostQueryParams) (*model.PostM, error) {
	var post model.PostM
	query.Status = known.POST_STATUS_PUBLIC
	db, err := buildWhere(u.db, query)
	if err != nil {
		log.C(ctx).Errorw("Failed to Get post from build where", "err", err)
		return nil, err
	}
	if err := db.First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Update 更新一条 post 数据库记录.
func (u *posts) Update(ctx context.Context, post *model.PostM) error {
	return u.db.Save(post).Error
}

// List 根据 offset 和 limit 返回指定用户的 post 列表.
func (u *posts) List(ctx context.Context, r *v1.ListPostRequest) (count int64, ret []*model.PostM, err error) {
	// 声明一个空的 []interface{} 切片用于存放查询条件
	queryWhere := make([]interface{}, 0)
	// 逐个创建查询条件并追加到 queryWhere 切片中
	queryWhere = append(queryWhere, []interface{}{"project_id", ctx.Value(known.XPrjectIDKey).(string)})
	queryWhere = append(queryWhere, []interface{}{"status", known.POST_STATUS_PUBLIC})
	if !gconvert.IsEmpty(r.CategoryID) {
		queryWhere = append(queryWhere, []interface{}{"category_id", r.CategoryID})
	}
	if !gconvert.IsEmpty(r.ColumnID) {
		queryWhere = append(queryWhere, []interface{}{"column_id", r.ColumnID})
	}
	if !gconvert.IsEmpty(r.PostID) {
		queryWhere = append(queryWhere, []interface{}{"post_id", "!=", r.PostID})
	}
	if !gconvert.IsEmpty(r.Query) {
		queryWhere = append(queryWhere, []interface{}{"title", "like", "%" + r.Query + "%"})
	}
	if !gconvert.IsEmpty(r.TagID) {
		queryWhere = append(queryWhere, []interface{}{"tag", "like", r.TagID + "%"})
	}
	if !gconvert.IsEmpty(r.IsTop) {
		queryWhere = append(queryWhere, []interface{}{"is_top", known.STATUS_OK})
	} else {
		queryWhere = append(queryWhere, []interface{}{"is_top", known.STATUS_DISABLE})
	}
	if !gconvert.IsEmpty(r.Type) {
		queryWhere = append(queryWhere, []interface{}{"type", "in", strings.Split(r.Type, ",")})
	} else {
		queryWhere = append(queryWhere, []interface{}{"type", known.POST_TYPE_POST})
	}
	db, err := buildQueryList(u.db, queryWhere, "*", "id desc", r.Offset, r.Limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list post from build where", "err", err)
		return
	}
	// Offset(-1).Limit(-1).Count(&count) 不能少。否则会count=0
	err = db.Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// Delete 根据 username, post_id 删除数据库 post 记录.
func (u *posts) Delete(ctx context.Context, username string, post_ids []string) error {
	err := u.db.Where("username = ? and post_id in (?)", username, post_ids).Delete(&model.PostM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
