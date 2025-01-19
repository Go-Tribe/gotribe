// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package comment

import (
	"context"
	"errors"
	"github.com/dengmengmian/ghelper/gmarkdown"
	"gotribe/pkg/api/v1"
	"gotribe/pkg/ip"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
)

// CommentBiz defines functions used to handle example request.
type CommentBiz interface {
	Create(ctx context.Context, username string, r *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error)
	Reply(ctx context.Context, username string, r *v1.ReplyCommentRequest) (*v1.CreateCommentResponse, error)
	Update(ctx context.Context, username, commentID string, r *v1.UpdateCommentRequest) error
	Get(ctx context.Context, commentID string) (*v1.GetCommentResponse, error)
	List(ctx context.Context, r *v1.ListCommentRequest) (*v1.ListCommentResponse, error)
}

// The implementation of CommentBiz interface.
type commentBiz struct {
	ds store.IStore
}

// Make sure that commentBiz implements the CommentBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ CommentBiz = (*commentBiz)(nil)

func New(ds store.IStore) *commentBiz {
	return &commentBiz{ds: ds}
}

// Create is the implementation of the `Create` method in CommentBiz interface.
func (b *commentBiz) Create(ctx context.Context, username string, r *v1.CreateCommentRequest) (*v1.CreateCommentResponse, error) {
	var commentM model.CommentM
	_ = copier.Copy(&commentM, r)
	commentM.Type = known.CommentPublish
	commentM.HtmlContent = gmarkdown.MdToHTML(commentM.Content)
	commentM.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	commentM.IP = ctx.Value(known.XClientIPKey).(string)
	country, city, regionName := ip.GeoIP(ctx, commentM.IP)
	commentM.Country = country
	commentM.RegionName = regionName
	commentM.City = city
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		return nil, errno.ErrUserNotFound
	}
	commentM.UserID = userM.UserID
	if err := b.ds.Comments().Create(ctx, &commentM); err != nil {
		return nil, err
	}

	return &v1.CreateCommentResponse{CommentID: commentM.CommentID}, nil
}

func (b *commentBiz) Reply(ctx context.Context, username string, r *v1.ReplyCommentRequest) (*v1.CreateCommentResponse, error) {
	var commentM model.CommentM
	_ = copier.Copy(&commentM, r)
	commentM.Type = known.CommentReply
	commentM.HtmlContent = gmarkdown.MdToHTML(commentM.Content)
	commentM.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
	commentM.ParentID = *r.ParentID
	commentM.IP = ctx.Value(known.XClientIPKey).(string)
	commentM.IP = ctx.Value(known.XClientIPKey).(string)
	country, city, regionName := ip.GeoIP(ctx, commentM.IP)
	commentM.Country = country
	commentM.RegionName = regionName
	commentM.City = city
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		return nil, errno.ErrUserNotFound
	}
	commentM.UserID = userM.UserID
	if err := b.ds.Comments().Create(ctx, &commentM); err != nil {
		return nil, err
	}

	return &v1.CreateCommentResponse{CommentID: commentM.CommentID}, nil
}

// Get is the implementation of the `Get` method in CommentBiz interface.
func (b *commentBiz) Get(ctx context.Context, commentID string) (*v1.GetCommentResponse, error) {
	comment, err := b.ds.Comments().Get(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrCommentNotFound
		}

		return nil, err
	}

	var resp v1.GetCommentResponse
	_ = copier.Copy(&resp, comment)

	resp.CreatedAt = comment.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = comment.UpdatedAt.Format(known.TimeFormat)

	return &resp, nil
}

// Update is the implementation of the `Update` method in CommentBiz interface.
func (b *commentBiz) Update(ctx context.Context, username, commentID string, r *v1.UpdateCommentRequest) error {
	commentM, err := b.ds.Comments().Get(ctx, commentID)
	if err != nil {
		return err
	}
	userM, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
	if err != nil {
		return errno.ErrUserNotFound
	}
	if userM.UserID != commentM.UserID {
		return errno.ErrPermissionDenied
	}
	if r.Content != nil {
		commentM.Content = *r.Content
		commentM.HtmlContent = gmarkdown.MdToHTML(commentM.Content)
	}
	if err := b.ds.Comments().Update(ctx, commentM); err != nil {
		return err
	}

	return nil
}

// List is the implementation of the `List` method in CommentBiz interface.
// List is the implementation of the `List` method in CommentBiz interface.
func (b *commentBiz) List(ctx context.Context, r *v1.ListCommentRequest) (*v1.ListCommentResponse, error) {
	// Step 1: Query top-level comments with pagination
	count, topLevelComments, err := b.ds.Comments().List(ctx, r.ObjectID, r.Offset, r.Limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list top-level comments from storage", "err", err)
		return nil, err
	}

	// Step 2: Extract top-level comment IDs
	topLevelCommentIDs := make([]uint, 0, count)
	for _, comment := range topLevelComments {
		topLevelCommentIDs = append(topLevelCommentIDs, comment.ID)
	}

	// Step 3: Query all replies for the object
	allReplies, err := b.ds.Comments().ListReplies(ctx, r.ObjectID)
	if err != nil {
		log.C(ctx).Errorw("Failed to list all replies from storage", "err", err)
		return nil, err
	}

	// Step 4: Build a map of replies for each comment
	replyMap := make(map[int][]*model.CommentM)
	for _, reply := range allReplies {
		replyMap[reply.ParentID] = append(replyMap[reply.ParentID], reply)
	}

	// Step 5: Query user information for all comments and replies
	userIDs := make([]string, 0, count+int64(len(allReplies)))
	userMap := make(map[string]*model.UserM)

	for _, comment := range topLevelComments {
		userIDs = append(userIDs, comment.UserID)
	}

	for _, reply := range allReplies {
		userIDs = append(userIDs, reply.UserID)
	}

	users, err := b.ds.Users().ListInUserID(ctx, userIDs)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	for _, user := range users {
		userMap[user.UserID] = user
	}

	// Step 6: Construct the final list of comments with replies and user information
	comments := make([]*v1.CommentInfo, 0, len(topLevelComments))
	for _, comment := range topLevelComments {
		commentInfo := b.buildCommentInfo(ctx, comment, replyMap, userMap)
		comments = append(comments, commentInfo)
	}

	return &v1.ListCommentResponse{TotalCount: count, Comments: comments}, nil
}

// buildCommentInfo constructs a CommentInfo with replies and user information
func (b *commentBiz) buildCommentInfo(ctx context.Context, comment *model.CommentM, replyMap map[int][]*model.CommentM, userMap map[string]*model.UserM) *v1.CommentInfo {
	var commentInfo v1.CommentInfo
	_ = copier.Copy(&commentInfo, comment)

	commentInfo.CreatedAt = comment.CreatedAt.Format(known.TimeFormat)
	commentInfo.UpdatedAt = comment.UpdatedAt.Format(known.TimeFormat)

	// Add user information to commentInfo
	user, exists := userMap[comment.UserID]
	if exists {
		commentInfo.UserID = user.UserID
		commentInfo.Nickname = user.Nickname
		commentInfo.Avatar = user.AvatarURL
	}

	// Add replies to the commentInfo
	commentInfo.Replies = make([]*v1.CommentInfo, 0)
	for _, reply := range replyMap[int(comment.ID)] {
		replyInfo := b.buildCommentInfo(ctx, reply, replyMap, userMap)
		commentInfo.Replies = append(commentInfo.Replies, replyInfo)
	}

	return &commentInfo
}
