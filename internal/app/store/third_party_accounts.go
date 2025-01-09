// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package store

import (
	"context"
	"gotribe/pkg/api/v1"

	"gorm.io/gorm"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
)

// AccountStore 定义了 account 模块在 store 层所实现的方法.
type AccountStore interface {
	Create(ctx context.Context, account *model.ThirdPartyAccountsM) error
	Get(ctx context.Context, accountWhere v1.AccountWhere) (*model.ThirdPartyAccountsM, error)
	Update(ctx context.Context, account *model.ThirdPartyAccountsM) error
}

// AccountStore 接口的实现.
type accounts struct {
	db *gorm.DB
}

// 确保 accounts 实现了 AccountStore 接口.
var _ AccountStore = (*accounts)(nil)

func newAccounts(db *gorm.DB) *accounts {
	return &accounts{db}
}

// Create 插入一条 account 记录.
func (u *accounts) Create(ctx context.Context, account *model.ThirdPartyAccountsM) error {
	return u.db.WithContext(ctx).Create(&account).Error
}

// Get 根据用户名查询指定 account 的数据库记录.
func (u *accounts) Get(ctx context.Context, accountWhere v1.AccountWhere) (*model.ThirdPartyAccountsM, error) {
	var account model.ThirdPartyAccountsM
	db, err := buildWhere(u.db, accountWhere)
	if err != nil {
		log.C(ctx).Errorw("Failed to Get account from build where", "err", err)
		return nil, err
	}
	if err := db.WithContext(ctx).First(&account).Error; err != nil {
		log.C(ctx).Errorw("Failed to Get account from sql", "err", err)
		return nil, err
	}

	return &account, nil
}

// Update 更新一条 account 数据库记录.
func (u *accounts) Update(ctx context.Context, account *model.ThirdPartyAccountsM) error {
	log.C(ctx).Infow("account update", "account:", account)
	return u.db.WithContext(ctx).Save(account).Error
}
