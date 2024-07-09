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

// Delete 删除指定的示例.
func (ctrl *ExampleController) Delete(c *gin.Context) {
	log.C(c).Infow("Delete example function called")

	if err := ctrl.b.Examples().Delete(c, c.GetString(known.XUsernameKey), c.Param("exampleID")); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
