// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package post

import (
	"context"
	"errors"
	util "gotribe/pkg/amount"
	"strings"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/v1"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// PostBiz defines functions used to handle post request.
type PostBiz interface {
	Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error)
	Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error
	Delete(ctx context.Context, username, postID string) error
	DeleteCollection(ctx context.Context, username string, postIDs []string) error
	Get(ctx context.Context, postID string) (*v1.GetPostResponse, error)
	List(ctx context.Context, r *v1.ListPostRequest) (*v1.ListPostResponse, error)
}

// The implementation of PostBiz interface.
type postBiz struct {
	ds store.IStore
}

// Make sure that postBiz implements the PostBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ PostBiz = (*postBiz)(nil)

func New(ds store.IStore) *postBiz {
	return &postBiz{ds: ds}
}

// Create is the implementation of the `Create` method in PostBiz interface.
func (b *postBiz) Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, r)
	postM.Author = username
	postM.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	if err := b.ds.Posts().Create(ctx, &postM); err != nil {
		return nil, err
	}

	return &v1.CreatePostResponse{PostID: postM.PostID}, nil
}

// Delete is the implementation of the `Delete` method in PostBiz interface.
func (b *postBiz) Delete(ctx context.Context, username, postID string) error {
	if err := b.ds.Posts().Delete(ctx, username, []string{postID}); err != nil {
		return err
	}

	return nil
}

// DeleteCollection is the implementation of the `DeleteCollection` method in PostBiz interface.
func (b *postBiz) DeleteCollection(ctx context.Context, username string, postIDs []string) error {
	if err := b.ds.Posts().Delete(ctx, username, postIDs); err != nil {
		return err
	}

	return nil
}

// Get is the implementation of the `Get` method in PostBiz interface.
func (b *postBiz) Get(ctx context.Context, postID string) (*v1.GetPostResponse, error) {
	post, err := b.ds.Posts().Get(ctx, v1.PostQueryParams{PostID: postID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPostNotFound
		}

		return nil, err
	}

	var resp v1.GetPostResponse
	_ = copier.Copy(&resp, post)
	// tag信息
	tagsM, err := b.ds.Tags().GetTags(ctx, strings.Split(post.Tag, ","))
	var tags []*v1.TagInfo
	_ = copier.Copy(&tags, tagsM)
	resp.Tags = tags

	// 分类信息
	categoryM, err := b.ds.Categories().Get(ctx, post.CategoryID)
	if err != nil {
		log.C(ctx).Errorw("Failed to get category from storage", "err", err)
		return nil, err
	}
	var category v1.CategoryInfo
	_ = copier.Copy(&category, categoryM)
	resp.Category = category
	resp.Images = strings.Split(post.Images, ",")
	resp.CreatedAt = post.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = post.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// Update is the implementation of the `Update` method in PostBiz interface.
func (b *postBiz) Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error {
	postM, err := b.ds.Posts().Get(ctx, v1.PostQueryParams{PostID: postID, Author: username})
	if err != nil {
		return err
	}

	if r.Title != nil {
		postM.Title = *r.Title
	}

	if r.Content != nil {
		postM.Content = *r.Content
	}

	if err := b.ds.Posts().Update(ctx, postM); err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `List` method in PostBiz interface.
func (b *postBiz) List(ctx context.Context, r *v1.ListPostRequest) (*v1.ListPostResponse, error) {
	count, list, err := b.ds.Posts().List(ctx, r)
	if err != nil {
		log.C(ctx).Errorw("Failed to list posts from storage", "err", err)
		return nil, err
	}

	posts := make([]*v1.PostInfo, 0, len(list))
	for _, item := range list {
		post := item
		// 分类信息
		categoryM, err := b.ds.Categories().Get(ctx, post.CategoryID)
		if err != nil {
			log.C(ctx).Errorw("Failed to get category from storage", "err", err)
			return nil, err
		}
		var category v1.CategoryInfo
		_ = copier.Copy(&category, categoryM)

		// 用户信息
		userM, err := b.ds.Users().Get(ctx, v1.UserWhere{UserID: post.UserID})
		if err != nil {
			log.C(ctx).Errorw("Failed to get user from storage", "err", err)
			return nil, err
		}

		// 标签信息
		tagSlice := strings.Split(post.Tag, ",")
		tagsM, err := b.ds.Tags().GetTags(ctx, tagSlice)
		if err != nil {
			log.C(ctx).Errorw("Failed to get tags from storage", "err", err)
			return nil, err
		}
		var tags []*v1.TagInfo
		_ = copier.Copy(&tags, tagsM)

		posts = append(posts, &v1.PostInfo{
			Author:      userM.Nickname,
			PostID:      post.PostID,
			Title:       post.Title,
			ColumnID:    post.ColumnID,
			Tag:         post.Tag,
			Icon:        post.Icon,
			Category:    category,
			Tags:        tags, // 添加标签信息
			Content:     post.Content,
			HtmlContent: post.HtmlContent,
			Location:    post.Location,
			People:      post.People,
			Time:        post.Time,
			Type:        post.Type,
			Video:       post.Video,
			UnitPrice:   util.FenToYuan(int(post.UnitPrice)),
			Images:      strings.Split(post.Images, ","),
			Description: post.Description,
			CreatedAt:   post.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:   post.UpdatedAt.Format(known.TimeFormat),
		})
	}

	return &v1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
