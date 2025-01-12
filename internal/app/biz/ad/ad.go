// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package ad

import (
	"context"
	"gotribe/pkg/api/v1"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
)

// AdBiz defines functions used to handle ad request.
type AdBiz interface {
	List(ctx context.Context, username string, offset, limit int) (*v1.ListAdResponse, error)
}

// The implementation of AdBiz interface.
type adBiz struct {
	ds store.IStore
}

// Make sure that adBiz implements the AdBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ AdBiz = (*adBiz)(nil)

func New(ds store.IStore) *adBiz {
	return &adBiz{ds: ds}
}

// List is the implementation of the `List` method in AdBiz interface.
func (b *adBiz) List(ctx context.Context, sceneID string, offset, limit int) (*v1.ListAdResponse, error) {
	count, list, err := b.ds.Ad().List(ctx, sceneID, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list ads from storage", "err", err)
		return nil, err
	}

	ads := make([]*v1.AdInfo, 0, len(list))
	for _, item := range list {
		ad := item
		ads = append(ads, &v1.AdInfo{
			AdID:        ad.AdID,
			URL:         ad.URL,
			URLType:     ad.URLType,
			Image:       ad.Image,
			Video:       ad.Video,
			Title:       ad.Title,
			Description: ad.Description,
			Sort:        ad.Sort,
			SceneID:     ad.SceneID,
			Ext:         ad.Ext,
			CreatedAt:   ad.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:   ad.UpdatedAt.Format(known.TimeFormat),
		})
	}

	return &v1.ListAdResponse{TotalCount: count, Ads: ads}, nil
}
