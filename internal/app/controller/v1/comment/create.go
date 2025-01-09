// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package comment

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/api/v1"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// Create 创建一条评论
func (ctrl *CommentController) Create(c *gin.Context) {
	log.C(c).Infow("Create example function called")

	var r v1.CreateCommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}
	resp, err := ctrl.b.Comments().Create(c, c.GetString(known.XUsernameKey), &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
