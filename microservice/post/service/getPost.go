package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) GetPost(ctx context.Context, req *pb.Request, resp *pb.Post) error {
	logger.Info("PostService GetPost")

	post, err := s.Dao.GetPostInfo(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if post == nil {
		return errno.ServerErr(errno.ErrItemNotExist, "")
	}

	comments, err := s.Dao.ListCommentByPostId(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	for _, comment := range comments {
		n, err := s.Dao.GetLikedNum(dao.Item{
			Id:     comment.Id,
			TypeId: 2,
		})
		if err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
		comment.LikeNum = uint32(n)
	}

	likeNum, err := s.Dao.GetLikedNum(dao.Item{
		Id:     req.Id,
		TypeId: 1,
	})
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp.LikeNum = post.LikeNum
	if likeNum != 0 {
		resp.LikeNum = uint32(likeNum)
	}
	resp.Id = post.Id
	resp.Content = post.Content
	resp.Title = post.Title
	resp.Time = post.LastEditTime
	resp.Category = post.Category
	resp.CreatorId = post.CreatorId
	resp.CreatorAvatar = post.CreatorAvatar
	resp.CreatorName = post.CreatorName
	resp.Comments = comments
	resp.CommentNum = uint32(len(comments))

	return nil
}
