// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/upload"
)

// Get 获取一个用户的详细信息.
func (ctrl *UserController) UploadResources(c *gin.Context) {
	log.C(c).Infow("post upload file", "user_id", c.GetString("user_id"), "request_id", c.GetString("request_id"))

	fileHeader, err := c.FormFile("file")
	if err != nil {
		core.WriteResponse(c, fmt.Errorf("failed to get form file: %w", err), nil)
		return
	}
	if fileHeader == nil {
		core.WriteResponse(c, errors.New("no file uploaded"), nil)
		return
	}
	if fileHeader.Size > known.DEFAULT_UPLOAD_SIZE {
		core.WriteResponse(c, fmt.Errorf("file size exceeds limit: %d bytes", known.DEFAULT_UPLOAD_SIZE), nil)
		return
	}

	provider := viper.GetString("upload-file.provider")
	if provider == "" {
		if viper.GetBool("enable-oss") {
			provider = "oss"
		} else {
			provider = "qiniu"
		}
	}
	uploadFile, err := upload.NewService(&upload.Options{
		Provider:  upload.Provider(provider),
		Endpoint:  viper.GetString("upload-file.endpoint"),
		AccessKey: viper.GetString("upload-file.access-key"),
		SecretKey: viper.GetString("upload-file.secret-key"),
		Bucket:    viper.GetString("upload-file.bucket"),
		Region:    viper.GetString("upload-file.region"),
	})
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	fileRes, err := uploadFile.UploadFile(fileHeader)
	if err != nil {
		core.WriteResponse(c, fmt.Errorf("failed to create upload file: %w", err), nil)
		return
	}
	fileRes.Domain = viper.GetString("cdn-domain")
	core.WriteResponse(c, nil, fileRes)
}
