// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package userEvent

import (
	"context"
	"github.com/jinzhu/copier"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/v1"
)

// UserEventBiz defines functions used to handle userEvent request.
type UserEventBiz interface {
	Create(ctx context.Context, r *v1.CreateUserEventRequest) error
}

// The implementation of UserEventBiz interface.
type userEventBiz struct {
	ds store.IStore
}

// Make sure that userEventBiz implements the UserEventBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ UserEventBiz = (*userEventBiz)(nil)

func New(ds store.IStore) *userEventBiz {
	return &userEventBiz{ds: ds}
}

func (b *userEventBiz) Create(ctx context.Context, r *v1.CreateUserEventRequest) error {
	var userEventM model.UserEventM
	if ctx.Value(known.XUsernameKey).(string) != "" {
		userInfo, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: ctx.Value(known.XUsernameKey).(string)})
		if err != nil {
			return err
		}
		userEventM.UserID = userInfo.UserID
	}
	_ = copier.Copy(&userEventM, r)
	if _, err := b.ds.UserEvents().Create(ctx, &userEventM); err != nil {
		return err
	}

	return nil
}
