package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
	"go.uber.org/zap"
)

func (s *PostService) ListSubPost(_ context.Context, req *pb.ListSubPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListSubPost")

	filter := &dao.PostModel{
		TypeName:   req.TypeName,
		MainPostId: req.MainPostId,
	}

	posts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = make([]*pb.Post, len(posts))
	for i, post := range posts {
		comments, err := s.Dao.ListCommentByPostId(post.Id)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
		}

		commentNum := s.Dao.GetCommentNumByPostId(post.Id)

		item := dao.Item{
			Id:       post.Id,
			TypeName: constvar.Post,
		}

		isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrRedis))
		}

		likeNum, err := s.Dao.GetLikedNum(item)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrRedis))
		}

		isFavorite, err := s.Dao.IsUserFavoritePost(req.UserId, post.Id)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
		}

		tags, err := s.Dao.ListTagsByPostId(post.Id)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
		}

		resp.List[i] = &pb.Post{
			Id:            post.Id,
			Title:         post.Title,
			Time:          util.FormatString(post.LastEditTime),
			CategoryId:    post.CategoryId,
			CreatorId:     post.CreatorId,
			CreatorName:   post.CreatorName,
			CreatorAvatar: post.CreatorAvatar,
			Content:       post.Content,
			CommentNum:    commentNum,
			LikeNum:       uint32(likeNum),
			IsLiked:       isLiked,
			IsFavorite:    isFavorite,
			Comments:      comments,
			Tags:          tags,
		}
	}

	return nil
}
