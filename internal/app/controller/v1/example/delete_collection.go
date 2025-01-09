// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package example

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// DeleteCollection 批量删除示例.
func (ctrl *ExampleController) DeleteCollection(c *gin.Context) {
	log.C(c).Infow("Batch delete comment function called")

	exampleIDs := c.QueryArray("exampleID")
	if err := ctrl.b.Examples().DeleteCollection(c, c.GetString(known.XUsernameKey), exampleIDs); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
