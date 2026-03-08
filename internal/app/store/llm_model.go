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

// LLMModelStore 定义 llm_model 在 store 层的方法.
type LLMModelStore interface {
	// GetByName 按 name 查询，多条时取第一条.
	GetByName(ctx context.Context, name string) (*model.LLMModelM, error)
}

type llmModels struct {
	db *gorm.DB
}

var _ LLMModelStore = (*llmModels)(nil)

func newLLMModels(db *gorm.DB) *llmModels {
	return &llmModels{db: db}
}

func (s *llmModels) GetByName(ctx context.Context, name string) (*model.LLMModelM, error) {
	var m model.LLMModelM
	err := s.db.WithContext(ctx).Where("name = ? AND status = 1", name).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
