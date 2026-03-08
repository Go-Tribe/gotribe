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

// ConversationLogStore 定义 conversation_log 在 store 层的方法.
type ConversationLogStore interface {
	Create(ctx context.Context, log *model.ConversationLogM) error
	Update(ctx context.Context, log *model.ConversationLogM) error
}

type conversationLogs struct {
	db *gorm.DB
}

var _ ConversationLogStore = (*conversationLogs)(nil)

func newConversationLogs(db *gorm.DB) *conversationLogs {
	return &conversationLogs{db: db}
}

func (s *conversationLogs) Create(ctx context.Context, log *model.ConversationLogM) error {
	return s.db.WithContext(ctx).Create(log).Error
}

func (s *conversationLogs) Update(ctx context.Context, log *model.ConversationLogM) error {
	return s.db.WithContext(ctx).Save(log).Error
}
