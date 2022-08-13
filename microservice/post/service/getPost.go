package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"go.uber.org/zap"
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

	for _, comment := range comments {
		item := dao.Item{
			Id:       comment.Id,
			TypeName: constvar.Comment,
		}

		n, err := s.Dao.GetLikedNum(item)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrRedis))
		}
		comment.LikeNum = uint32(n)

		isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrRedis))
		}
		comment.IsLiked = isLiked
	}

	item := dao.Item{
		Id:       req.Id,
		TypeName: constvar.Post,
	}

	likeNum, err := s.Dao.GetLikedNum(item)
	if err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrRedis))
	}
	resp.LikeNum = post.LikeNum
	if likeNum != 0 {
		resp.LikeNum = uint32(likeNum)
	}

	isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
	if err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrRedis))
	}
	resp.IsLiked = isLiked

	isFavorite, err := s.Dao.IsUserFavoritePost(req.UserId, req.Id)
	if err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
	}
	resp.IsFavorite = isFavorite

	tags, err := s.Dao.ListTagsByPostId(post.Id)
	if err != nil {
		logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
	}

	resp.Id = post.Id
	resp.Content = post.Content
	resp.Title = post.Title
	resp.Time = post.LastEditTime
	resp.CategoryId = post.CategoryId
	resp.CreatorId = post.CreatorId
	resp.CreatorAvatar = post.CreatorAvatar
	resp.CreatorName = post.CreatorName
	resp.Comments = comments
	resp.CommentNum = uint32(len(comments))
	resp.Tags = tags

	return nil
}
