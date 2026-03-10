// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"time"

	"gotribe/internal/pkg/model"

	"gorm.io/gorm"
)

// VerificationCodeStore 定义 verification_code 在 store 层的方法.
type VerificationCodeStore interface {
	Create(ctx context.Context, v *model.VerificationCodeM) error
	// GetLatestUnverified 取指定 target+trigger 下未过期、未验证的最新一条（用于校验验证码）.
	GetLatestUnverified(ctx context.Context, target, trigger string) (*model.VerificationCodeM, error)
	Update(ctx context.Context, v *model.VerificationCodeM) error
}

type verificationCodes struct {
	db *gorm.DB
}

var _ VerificationCodeStore = (*verificationCodes)(nil)

func newVerificationCodes(db *gorm.DB) *verificationCodes {
	return &verificationCodes{db: db}
}

func (s *verificationCodes) Create(ctx context.Context, v *model.VerificationCodeM) error {
	return s.db.WithContext(ctx).Create(v).Error
}

func (s *verificationCodes) GetLatestUnverified(ctx context.Context, target, trigger string) (*model.VerificationCodeM, error) {
	var m model.VerificationCodeM
	now := time.Now()
	err := s.db.WithContext(ctx).
		Where("target = ? AND trigger = ? AND verified = 0 AND expire_at > ?", target, trigger, now).
		Order("id DESC").
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *verificationCodes) Update(ctx context.Context, v *model.VerificationCodeM) error {
	return s.db.WithContext(ctx).Save(v).Error
}
