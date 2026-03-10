// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	"gotribe/pkg/api/v1"
	"gotribe/pkg/auth"
	"gotribe/pkg/email"
	"gotribe/pkg/token"
	"regexp"
	"sync"
	"time"
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
	SendVerificationCode(ctx context.Context, r *v1.SendVerificationCodeRequest, opts *email.Options, expireMinutes int) error
	Register(ctx context.Context, r *v1.RegisterRequest) error
	VerificationCodeLogin(ctx context.Context, r *v1.VerificationCodeLoginRequest) (*v1.LoginResponse, error)
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
	if _, err := b.ds.Users().Create(ctx, &userM); err != nil {
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
	if !gconvert.IsEmpty(user.Birthday) {
		resp.Birthday = user.Birthday.Format(known.TimeFormatShort)
	}
	// 获取用户积分
	point, err := b.ds.PointAvailable().SumPoints(ctx, user.UserID, ctx.Value(known.XProjectIDKey).(string))
	if err != nil {
		return nil, err
	}
	if point != nil {
		resp.Point = *point
	} else {
		resp.Point = 0 // 或者其他默认值
	}
	resp.CreatedAt = user.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = user.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// List 是 UserBiz 接口中 `List` 方法的实现.
func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit, v1.UserWhere{ProjectID: ctx.Value(known.XProjectIDKey).(string)})
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
	// 获取用户信息
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrUserNotFound
		}
		log.C(ctx).Errorw("Failed to get user", "username", username, "err", err)
		return err
	}

	// 更新字段
	if user.Email != nil {
		userM.Email = *user.Email
	}
	if user.Nickname != nil {
		userM.Nickname = *user.Nickname
	}
	if user.Phone != nil {
		userM.Phone = *user.Phone
	}
	if user.AvatarURL != nil {
		userM.AvatarURL = *user.AvatarURL
	}
	if user.Sex != nil {
		userM.Sex = *user.Sex
	}
	if user.Background != nil {
		userM.Background = *user.Background
	}
	if user.Ext != nil {
		userM.Ext = *user.Ext
	}

	// 更新生日字段
	if user.Birthday != nil {
		birthday, err := time.Parse(known.TimeFormatShort, *user.Birthday)
		if err != nil {
			log.C(ctx).Warnw("Failed to parse birthday, skipping update", "username", username, "err", err)
		} else {
			userM.Birthday = &birthday
		}
	}

	// 使用事务确保原子性
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		log.C(ctx).Errorw("Failed to update user", "username", username, "err", err)
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
			return b.createUserAndAccount(ctx, openID, known.LOGIN_TYPE_WXMINI)
		} else {
			return nil, err // 其他错误直接返回。
		}
	}
	userInfo, err := b.ds.Users().Get(ctx, v1.UserWhere{UserID: accountM.UserID})
	if err != nil {
		return nil, err
	}
	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(userInfo.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t, Username: userInfo.Username}, nil
}

// genVerificationCode 生成 6 位数字验证码.
func genVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SendVerificationCode 发送验证码：落库并发邮件.
func (b *userBiz) SendVerificationCode(ctx context.Context, r *v1.SendVerificationCodeRequest, opts *email.Options, expireMinutes int) error {
	code, err := genVerificationCode()
	if err != nil {
		return errno.InternalServerError
	}
	if expireMinutes <= 0 {
		expireMinutes = 10
	}
	expireAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)
	projectID := ""
	if v := ctx.Value(known.XProjectIDKey); v != nil {
		if s, ok := v.(string); ok {
			projectID = s
		}
	}
	vc := &model.VerificationCodeM{
		Channel:   "email",
		Trigger:   r.Trigger,
		Target:    r.Email,
		Code:      code,
		ProjectID: projectID,
		ExpireAt:  expireAt,
		Verified:  0,
	}
	if err := b.ds.VerificationCodes().Create(ctx, vc); err != nil {
		return err
	}
	if opts != nil {
		body := email.BuildVerificationMailBody(code, expireMinutes)
		if err := email.Send(opts, r.Email, email.DefaultSubject, body); err != nil {
			log.C(ctx).Errorw("Send verification email failed", "email", r.Email, "err", err)
			return errno.ErrEmailSendFailed
		}
	}
	return nil
}

// Register 校验邮箱验证码后注册用户.
func (b *userBiz) Register(ctx context.Context, r *v1.RegisterRequest) error {
	rec, err := b.ds.VerificationCodes().GetLatestUnverified(ctx, r.Email, "register")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrVerificationCodeInvalid
		}
		return err
	}
	if rec.Code != r.Code {
		return errno.ErrVerificationCodeInvalid
	}
	rec.Verified = 1
	if err := b.ds.VerificationCodes().Update(ctx, rec); err != nil {
		return err
	}
	createReq := &v1.CreateUserRequest{
		Username: r.Username,
		Password: r.Password,
		Nickname: r.Nickname,
		Email:    r.Email,
		Phone:    r.Phone,
	}
	return b.Create(ctx, createReq)
}

// VerificationCodeLogin 验证码登录：先校验验证码（trigger=login），有账号则登录，无则自动注册再登录. Target 目前仅支持邮箱.
func (b *userBiz) VerificationCodeLogin(ctx context.Context, r *v1.VerificationCodeLoginRequest) (*v1.LoginResponse, error) {
	rec, err := b.ds.VerificationCodes().GetLatestUnverified(ctx, r.Target, "login")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrVerificationCodeInvalid
		}
		return nil, err
	}
	if rec.Code != r.Code {
		return nil, errno.ErrVerificationCodeInvalid
	}
	rec.Verified = 1
	if err := b.ds.VerificationCodes().Update(ctx, rec); err != nil {
		return nil, err
	}

	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Email: r.Target})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if userM != nil {
		t, err := token.Sign(userM.Username)
		if err != nil {
			return nil, errno.ErrSignToken
		}
		return &v1.LoginResponse{Token: t, Username: userM.Username}, nil
	}

	// 无账号：自动注册再登录（仅支持邮箱，生成用户名与随机密码）
	username, err1 := genRandomPassword(15)
	nickname := "用户" + codeSuffix(4)
	passRaw, err := genRandomPassword(8)
	if err != nil || err1 != nil {
		return nil, errno.InternalServerError
	}
	createReq := &v1.CreateUserRequest{
		Username: username,
		Password: passRaw,
		Nickname: nickname,
		Email:    r.Target,
		Phone:    "",
	}
	if err := b.Create(ctx, createReq); err != nil {
		if err == errno.ErrUserAlreadyExist {
			// 用户名冲突，重试一次加后缀
			username = username + codeSuffix(4)
			createReq.Username = username
			if err := b.Create(ctx, createReq); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	t, err := token.Sign(username)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t, Username: username}, nil
}

// emailToUsername 将邮箱转为合法用户名（仅保留字母数字下划线）.
func emailToUsername(email string) string {
	var b []byte
	for _, c := range email {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			b = append(b, byte(c))
		} else if c == '@' || c == '.' {
			b = append(b, '_')
		}
	}
	s := string(b)
	if len(s) > 18 {
		s = s[:18]
	}
	if len(s) < 6 {
		s = s + codeSuffix(6-len(s))
	}
	return s
}

func codeSuffix(n int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	x, _ := rand.Int(rand.Reader, max)
	s := fmt.Sprintf("%0*d", n, x.Int64())
	return s[:n]
}

func genRandomPassword(n int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		x, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[x.Int64()]
	}
	return string(b), nil
}
