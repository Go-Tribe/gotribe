// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package errno

var ErrCommentNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.CommentNotFound", Message: "评论不存在"}
var ErrPermissionDenied = &Errno{HTTP: 403, Code: "AccessDenied.PermissionDenied", Message: "权限不足"}
