// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package app_version

import (
	"gotribe/internal/app/biz/app_version"
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// LatestRelease 获取当前产品+平台的最新发布版本.
// Header: x-product, x-platform, x-platform-version-code, x-platform-version-name.
// 除基础版本字段外返回 needForceUpdate：最新版 ForceUpdate=1 或 客户端版本号 < 最新版最低兼容版本号 时为 true.
func (ctrl *AppVersionController) LatestRelease(c *gin.Context) {
	log.C(c).Infow("Get latest release called")

	productName := c.GetHeader(known.XProductKey)
	platform := c.GetHeader(known.XPlatformKey)
	if productName == "" || platform == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}

	codeStr := c.GetHeader(known.XPlatformVersionCodeKey)
	clientVersionCode := app_version.ParseClientVersionCode(codeStr)

	resp, err := ctrl.b.AppVersions().GetLatestRelease(c, productName, platform, clientVersionCode)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
