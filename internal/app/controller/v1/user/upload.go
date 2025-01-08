// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"github.com/spf13/viper"
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/upload"

	"github.com/gin-gonic/gin"
)

// Get 获取一个用户的详细信息.
func (ctrl *UserController) UploadResources(c *gin.Context) {
	log.C(c).Infow("post upload file")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	if fileHeader.Size > known.DEFAULT_UPLOAD_SIZE {
		core.WriteResponse(c, err, nil)
		return
	}

	uploadFile, err := upload.NewUploadFile(
		viper.GetString("upload-file.endpoint"),
		viper.GetString("upload-file.access-key"),
		viper.GetString("upload-file.secret-key"),
		viper.GetString("upload-file.bucket"),
		viper.GetBool("enable-oss"),
	)
	fileRes, err := uploadFile.UploadFile(fileHeader)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	fileRes.Domain = viper.GetString("cdn-domain")

	core.WriteResponse(c, nil, fileRes)
}
