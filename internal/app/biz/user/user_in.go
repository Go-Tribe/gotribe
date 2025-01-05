package user

import (
	"context"
	"fmt"
	"github.com/dengmengmian/ghelper/gid"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/model"
	v1 "gotribe/pkg/api/v1"
	"gotribe/pkg/token"
	"math/rand"
)

func (b *userBiz) createUserAndAccount(ctx context.Context, openID, platformType string) (*v1.LoginResponse, error) {
	userid := gid.GenShortID()
	// 先建用户，再新增拓展用户信息
	userM := model.UserM{
		UserID:    userid,
		ProjectID: ctx.Value(known.XPrjectIDKey).(string),
		Username:  generateRandomUserName(),
	}
	err := b.ds.Users().Create(ctx, &userM)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accountM := model.ThirdPartyAccountsM{
		UserID:   userid,
		OpenID:   openID,
		Platform: platformType,
	}
	err = b.ds.ThirdPartyAccounts().Create(ctx, &accountM)
	if err != nil {
		return nil, fmt.Errorf("failed to create third party account: %w", err)
	}

	t, err := token.Sign(userid)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t}, nil
}

func generateRandomUserName() string {
	// 定义用户名前缀
	prefix := "wxmini_"

	// 生成随机数字部分
	randomNumber := rand.Intn(999999) // 生成一个0到999999之间的随机数

	// 将随机数字转换为字符串
	randomNumberStr := fmt.Sprintf("%06d", randomNumber) // 补足6位数字，不足的前面补0

	// 拼接用户名
	userName := prefix + randomNumberStr

	return userName
}
