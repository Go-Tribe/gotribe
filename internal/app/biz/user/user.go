// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"context"
	"errors"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	"gotribe/pkg/api/v1"
	"gotribe/pkg/auth"
	"gotribe/pkg/token"
	"regexp"
	"sync"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// UserBiz 定义了 user 模块在 biz 层所实现的方法.
type UserBiz interface {
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
	Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
	DeleteToTx(ctx context.Context, username string) error
	WxMiniLogin(ctx context.Context, r *v1.WechatMiniLoginRequest, openID string) (*v1.LoginResponse, error)
}

// UserBiz 接口的实现.
type userBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

// New 创建一个实现了 UserBiz 接口的实例.
func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// ChangePassword 是 UserBiz 接口中 `ChangePassword` 方法的实现.
func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		return err
	}

	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

// Login 是 UserBiz 接口中 `Login` 方法的实现.
func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取登录用户的所有信息
	user, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: r.Username})
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

// Create 是 UserBiz 接口中 `Create` 方法的实现.
func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)
	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

// Get 是 UserBiz 接口中 `Get` 方法的实现.
func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)

	resp.CreatedAt = user.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = user.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// List 是 UserBiz 接口中 `List` 方法的实现.
func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit, v1.UserWhere{ProjectID: ctx.Value(known.XPrjectIDKey).(string)})
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	// 使用 goroutine 提高接口性能
	for _, item := range list {
		user := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				// count, _, err := b.ds.Posts().List(ctx, user.Username, 0, 0)
				if err != nil {
					log.C(ctx).Errorw("Failed to list posts", "err", err)
					return err
				}

				m.Store(user.ID, &v1.UserInfo{
					UserID:   user.UserID,
					Username: user.Username,
					Nickname: user.Nickname,
					Email:    user.Email,
					Phone:    user.Email,
					// PostCount: count,
					CreatedAt: user.CreatedAt.Format(known.TimeFormat),
					UpdatedAt: user.UpdatedAt.Format(known.TimeFormat),
				})

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.UserInfo))
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// Update 是 UserBiz 接口中 `Update` 方法的实现.
func (b *userBiz) Update(ctx context.Context, username string, user *v1.UpdateUserRequest) error {
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		return err
	}

	if user.Email != nil {
		userM.Email = *user.Email
	}

	if user.Nickname != nil {
		userM.Nickname = *user.Nickname
	}

	if user.Phone != nil {
		userM.Phone = *user.Phone
	}

	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

// Delete 是 UserBiz 接口中 `Delete` 方法的实现.
func (b *userBiz) Delete(ctx context.Context, username string) error {
	if err := b.ds.Users().Delete(ctx, username); err != nil {
		return err
	}

	return nil
}

// DeleteToTx 是 UserBiz 接口中 `Delete` 事务方法的实现.
func (b *userBiz) DeleteToTx(ctx context.Context, username string) error {
	err := b.ds.TX(ctx, func(ctx context.Context) error {
		if err := b.ds.Users().Delete(ctx, username); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *userBiz) WxMiniLogin(ctx context.Context, r *v1.WechatMiniLoginRequest, openID string) (*v1.LoginResponse, error) {
	log.C(ctx).Infow("WxMiniLogin function called")
	if openID == "" {
		return nil, errno.ErrUserNotFound
	}

	// 获取拓展用户信息
	accountM, err := b.ds.ThirdPartyAccounts().Get(ctx, v1.AccountWhere{OpenID: openID})
	log.C(ctx).Infow("get account info", accountM, "err:", err)
	if err != nil {
		// 如果不存在则新建
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return b.createUserAndAccount(ctx, openID, r.Type)
		} else {
			return nil, err // 其他错误直接返回。
		}
	}
	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(accountM.UserID)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t}, nil
}
