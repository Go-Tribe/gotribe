// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type CommentM struct {
	gorm.Model
	CommentID   string `gorm:"type:char(10);uniqueIndex;example:唯一字符ID/分布式ID" json:"commentID"`
	ProjectID   string `gorm:"type:char(10);not null;index;example:项目ID;" json:"projectID"`
	Content     string `gorm:"not null;type:longtext;not null;example:内容" json:"content"`
	HtmlContent string `gorm:"not null;type:longtext;not null;example:HTML内容" json:"htmlContent"`
	Status      uint   `gorm:"type:tinyint;not null;index;default:1;example:状态，1-待审核；2-审核通过" json:"status,omitempty"`
	ObjectID    string `gorm:"type:char(10);not null;index;example:评论主题ID" json:"objectID"`
	ObjectType  uint   `gorm:"type:tinyint;not null;default:1;index;example:评论对象类型，1-文章；2-商品" json:"objectType"`
	Type        uint   `gorm:"type:tinyint;not null;default:1;example:评论类型，1-评论；2-回复" json:"type"`
	UserID      string `gorm:"type:char(10);not null;index;example:用户ID" json:"userID"`
	ToUserID    string `gorm:"type:char(10);not null;index;example:被评论用户ID" json:"toUserID"`
	PID         int    `gorm:"type:int;not null;default:0;example:父评论ID" json:"pid"`
	ReplyToID   int    `gorm:"type:int;not null;default:0;example:回复的评论ID" json:"ReplyToID"`
	Hot         int    `gorm:"type:int;default:0;example:热度" json:"hot"`
	Like        int    `gorm:"type:int;default:0;example:点赞数" json:"like"`
	Dislike     int    `gorm:"type:int;default:0;example:踩数" json:"dislike"`
}

func (c *CommentM) TableName() string {
	return "example"
}

func (c *CommentM) BeforeCreate(tx *gorm.DB) (err error) {
	// Encrypt the user password.
	c.CommentID = gid.GenShortID()
	return nil
}
