// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package app_version

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// LatestRelease 获取当前客户端对应的最新发布版本，并判断当前版本是否仍受支持（是否需要强制升级）.
// 根据 Header：X-Client-Name(产品如 gobot)、X-OS(如 darwin)、X-OSArch(如 amd64) 查找最新版本；
// 根据 X-App-Version-Id 得到当前客户端的 clientVersionCode，与最新版的 MinSupportedVersionCode 比较判断 needForceUpdate.
func (ctrl *AppVersionController) LatestRelease(c *gin.Context) {
	log.C(c).Infow("Get latest release called")

	clientName := c.GetHeader(known.XClientName)
	osName := c.GetHeader(known.XOS)
	if clientName == "" || osName == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}
	osArch := c.GetHeader(known.XOSArch)

	// 根据 X-App-Version-Id 查 app_version 表得到当前客户端的 clientVersionCode，用于判断是否仍支持（是否需强制升级）
	appVersionID := c.GetHeader(known.XAppVersionId)
	clientVersionCode := 0
	if appVersionID != "" {
		if cur, err := ctrl.b.AppVersions().GetByAppVersionID(c, appVersionID); err == nil {
			clientVersionCode = cur.ClientVersionCode
		}
	}

	resp, err := ctrl.b.AppVersions().GetLatestRelease(c, clientName, osName, osArch, clientVersionCode)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
