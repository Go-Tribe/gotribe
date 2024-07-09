// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// GetCategoryResponse 指定了 `GET /v1/category/{categoryID}` 接口的返回参数.
type GetCategoryResponse CategoryInfo

// CategoryInfo 指定了分类详细信息.
type CategoryInfo struct {
	CategoryID  string `json:"categoryID,omitempty"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
