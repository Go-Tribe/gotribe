// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package product

import (
	"context"
	"errors"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe/pkg/api/v1"
	"strings"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
)

// ProductBiz defines functions used to handle product request.
type ProductBiz interface {
	Get(ctx context.Context, productID string) (*v1.GetProductResponse, error)
	List(ctx context.Context, r *v1.ListProductRequest) (*v1.ListProductResponse, error)
}

// The implementation of ProductBiz interface.
type productBiz struct {
	ds store.IStore
}

// Make sure that productBiz implements the ProductBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ ProductBiz = (*productBiz)(nil)

func New(ds store.IStore) *productBiz {
	return &productBiz{ds: ds}
}

// Get is the implementation of the `Get` method in ProductBiz interface.
func (b *productBiz) Get(ctx context.Context, productID string) (*v1.GetProductResponse, error) {
	product, err := b.ds.Products().Get(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPageNotFound
		}
		return nil, err
	}

	var resp v1.GetProductResponse
	_ = copier.Copy(&resp, product)
	// 如果 tag 存在则取出 tag 信息
	if !gconvert.IsEmpty(product.Tag) {
		// tag信息
		tagsM, err := b.ds.Tags().GetTags(ctx, strings.Split(product.Tag, ","))
		if err != nil {
			log.C(ctx).Errorw("Failed to get tags from storage", "err", err)
			return nil, err
		}
		var tags []*v1.TagInfo
		_ = copier.Copy(&tags, tagsM)
		resp.Tags = tags
	}
	// sku信息
	skusM, err := b.ds.ProductSKUs().List(ctx, productID)
	if err != nil {
		log.C(ctx).Errorw("Failed to list product skus from storage", "err", err)
		return nil, err
	}
	var skus []*v1.ProductSKUInfo
	_ = copier.Copy(&skus, skusM)
	resp.Skus = skus
	resp.CreatedAt = product.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = product.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// List is the implementation of the `List` method in ProductBiz interface.
func (b *productBiz) List(ctx context.Context, r *v1.ListProductRequest) (*v1.ListProductResponse, error) {
	count, list, err := b.ds.Products().List(ctx, r.CategoryID, r.Offset, r.Limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list products from storage", "err", err)
		return nil, err
	}

	products := make([]*v1.ProductInfo, 0, len(list))
	for _, item := range list {
		product := item
		var tags []*v1.TagInfo
		if !gconvert.IsEmpty(item.Tag) {
			// tag信息
			tagsM, err := b.ds.Tags().GetTags(ctx, strings.Split(item.Tag, ","))
			if err != nil {
				log.C(ctx).Errorw("Failed to get tags from storage", "err", err)
				return nil, err
			}
			var tags []*v1.TagInfo
			_ = copier.Copy(&tags, tagsM)
		}
		// sku信息
		skusM, err := b.ds.ProductSKUs().List(ctx, item.ProductID)
		if err != nil {
			log.C(ctx).Errorw("Failed to list product skus from storage", "err", err)
			return nil, err
		}
		var skus []*v1.ProductSKUInfo
		_ = copier.Copy(&skus, skusM)
		products = append(products, &v1.ProductInfo{
			CategoryID:    product.CategoryID,
			Image:         product.Image,
			ProductNumber: product.ProductNumber,
			Description:   product.Description,
			HtmlContent:   product.HtmlContent,
			Video:         product.Video,
			BuyLimit:      product.BuyLimit,
			ProjectID:     product.ProjectID,
			ProductID:     product.ProductID,
			Title:         product.Title,
			Content:       product.Content,
			Tags:          tags,
			Skus:          skus,
			CreatedAt:     product.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:     product.UpdatedAt.Format(known.TimeFormat),
		})
	}

	return &v1.ListProductResponse{TotalCount: count, Products: products}, nil
}
