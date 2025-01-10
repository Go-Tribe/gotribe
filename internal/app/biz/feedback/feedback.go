// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package feedBack

import (
	"context"
	"errors"
	"gotribe/pkg/api/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"
)

// FeedBackBiz defines functions used to handle comment request.
type FeedBackBiz interface {
	Create(ctx context.Context, username string, r *v1.CreateFeedBackRequest) (*v1.CreateFeedBackResponse, error)
}

// The implementation of FeedBackBiz interface.
type feedBackBiz struct {
	ds store.IStore
}

// Make sure that feedBackBiz implements the FeedBackBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ FeedBackBiz = (*feedBackBiz)(nil)

func New(ds store.IStore) *feedBackBiz {
	return &feedBackBiz{ds: ds}
}

// Create is the implementation of the `Create` method in FeedBackBiz interface.
func (b *feedBackBiz) Create(ctx context.Context, username string, r *v1.CreateFeedBackRequest) (*v1.CreateFeedBackResponse, error) {
	var feedBackM model.FeedbackM
	_ = copier.Copy(&feedBackM, r)
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}
	feedBackM.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	feedBackM.UserID = userM.UserID

	if err := b.ds.Feedbacks().Create(ctx, &feedBackM); err != nil {
		return nil, err
	}

	return &v1.CreateFeedBackResponse{FeedBackID: feedBackM.ID}, nil
}
