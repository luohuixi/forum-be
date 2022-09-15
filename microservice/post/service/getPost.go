package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"
)

func (s *PostService) GetPost(_ context.Context, req *pb.Request, resp *pb.Post) error {
	logger.Info("PostService GetPost")

	post, err := s.Dao.GetPostInfo(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if post == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "post-"+strconv.Itoa(int(req.Id)))
	}

	comments, err := s.Dao.ListCommentByPostId(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	s.processComments(req.UserId, comments)

	resp.IsLiked, resp.IsCollection, resp.LikeNum, resp.Tags, resp.CommentNum = s.getPostInfo(post.Id, req.UserId)

	if resp.LikeNum == 0 {
		resp.LikeNum = post.LikeNum
	}

	resp.Id = post.Id
	resp.Content = post.Content
	resp.CompiledContent = post.CompiledContent
	resp.Title = post.Title
	resp.Time = post.LastEditTime
	resp.Category = post.Category
	resp.CreatorId = post.CreatorId
	resp.CreatorAvatar = post.CreatorAvatar
	resp.CreatorName = post.CreatorName
	resp.Comments = comments
	resp.ContentType = post.ContentType

	return nil
}
