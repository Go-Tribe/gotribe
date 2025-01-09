// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gotribe/internal/pkg/known"

	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// CommentStore 定义了 example 模块在 store 层所实现的方法.
type CommentStore interface {
	Create(ctx context.Context, comment *model.CommentM) error
	Get(ctx context.Context, commentID string) (*model.CommentM, error)
	Update(ctx context.Context, comment *model.CommentM) error
	List(ctx context.Context, objectID string, offset, limit int) (int64, []*model.CommentM, error)
	ListReplies(ctx context.Context, topLevelCommentIDs []string) ([]*model.CommentM, error)
}

// CommentStore 接口的实现.
type comments struct {
	db *gorm.DB
}

// 确保 comments 实现了 CommentStore 接口.
var _ CommentStore = (*comments)(nil)

func newComments(db *gorm.DB) *comments {
	return &comments{db}
}

// Create 插入一条记录.
func (u *comments) Create(ctx context.Context, comment *model.CommentM) error {
	return u.db.Create(&comment).Error
}

// Get 根据 commentID 查询数据库记录.
func (u *comments) Get(ctx context.Context, commentID string) (*model.CommentM, error) {
	var comment model.CommentM
	if err := u.db.Where("comment_id = ? and status = ?", commentID, known.AuditPass).First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

// Update 更新一条数据库记录.
func (u *comments) Update(ctx context.Context, comment *model.CommentM) error {
	return u.db.Save(comment).Error
}

// List 根据 offset 和 limit 返回列表.
func (u *comments) List(ctx context.Context, objectID string, offset, limit int) (count int64, ret []*model.CommentM, err error) {
	err = u.db.Where("object_id = ? and status = ? and pid = 0", objectID, known.AuditPass).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (c *comments) ListReplies(ctx context.Context, topLevelCommentIDs []string) ([]*model.CommentM, error) {
	var replies []*model.CommentM
	if err := c.db.WithContext(ctx).Where("pid IN ?", topLevelCommentIDs).Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}
