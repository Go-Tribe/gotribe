// Copyright 2023 Innkeeper dengmengmian(麻凡) <my@dengmengmian.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/dengmengmian/solocms

package point

import (
	"context"
	"fmt"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/v1"
)

type PointBiz interface {
	SubPoints(ctx context.Context, username, projectID, types, reason, eventID string, points float64) error
	GetAvailablePoints(ctx context.Context, userID, projectID string) (*float64, error)
}

type pointBiz struct {
	ds store.IStore
}

// 确保 goodsBiz 实现了 GoodsBiz 接口.
var _ PointBiz = (*pointBiz)(nil)

// New 创建一个实现了 GoodsBiz 接口的实例.
func New(ds store.IStore) *pointBiz {
	return &pointBiz{ds: ds}
}

// 获取用户可用积分
func (b *pointBiz) GetAvailablePoints(ctx context.Context, userID, projectID string) (*float64, error) {
	// 获取用户积分
	point, err := b.ds.PointAvailable().SumPoints(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available points: %v", err)
	}
	return point, nil
}

// 扣减积分
func (b *pointBiz) SubPoints(ctx context.Context, username, projectID, types, reason, eventID string, points float64) error {
	if points <= 0 {
		return fmt.Errorf("points must be positive")
	}

	err := b.ds.TX(ctx, func(ctx context.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.C(ctx).Errorw("panic occurred in transaction", "error", r)
			}
		}()

		// 查询用户 useriD
		userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
		if err != nil {
			return err
		}

		list, err := b.ds.PointAvailable().GetAll(ctx, userM.UserID, projectID)
		if err != nil {
			return err
		}

		if len(list) == 0 {
			return fmt.Errorf("no available points to deduct for user %s", username)
		}

		// 新增记录
		log.C(ctx).Infow("扣减积分", "points", points)
		logID, err := b.ds.PointLog().Create(ctx, &model.PointLogM{
			UserID:  userM.UserID,
			Points:  -points,
			Type:    types,
			Reason:  reason,
			EventID: eventID,
		})
		if err != nil {
			return err
		}

		// 循环扣减可用积分记录并在扣减积分表记录
		for _, item := range list {
			deductedPoints := pointmin(points, item.Points)
			points -= deductedPoints

			log.C(ctx).Infow("扣减积分", "deducted_points", deductedPoints)

			err = b.createPointDeduction(ctx, userM.UserID, int64(logID), int64(item.ID), deductedPoints)
			if err != nil {
				log.C(ctx).Infow("扣减积分失败", "points", deductedPoints)
				return err
			}

			item.Points -= deductedPoints
			if item.Points == 0 {
				item.Status = 2
			}

			err = b.updatePointAvailable(ctx, item)
			if err != nil {
				return err
			}

			if points <= 0 {
				break
			}
		}

		if points > 0 {
			return fmt.Errorf("not enough points to deduct, remaining points: %.2f", points)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *pointBiz) createPointDeduction(ctx context.Context, userID string, logID int64, availablePointsID int64, points float64) error {
	return b.ds.PointDeduction().Create(ctx, &model.PointDeductionM{
		UserID:            userID,
		Points:            points,
		PointsDetailID:    int(logID),
		AvailablePointsID: int(availablePointsID),
	})
}

func (b *pointBiz) updatePointAvailable(ctx context.Context, item *model.PointAvailableM) error {
	return b.ds.PointAvailable().Update(ctx, item)
}

func pointmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
