// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package user

import (
	"errors"
	"fmt"
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// UploadResources 上传用户资源（如图片）.
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

	uploadFile := upload.DefaultService()
	if uploadFile == nil {
		core.WriteResponse(c, errors.New("upload service not configured"), nil)
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
