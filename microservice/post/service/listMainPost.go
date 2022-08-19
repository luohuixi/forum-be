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

func (s *PostService) ListMainPost(_ context.Context, req *pb.ListMainPostRequest, resp *pb.ListPostResponse) error {
	logger.Info("PostService ListMainPost")

	filter := &dao.PostModel{
		TypeName:   req.TypeName,
		CategoryId: req.CategoryId,
	}

	posts, err := s.Dao.ListPost(filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.List = make([]*pb.Post, len(posts))
	for i, post := range posts {

		// limit max len of post content
		content := post.Content
		if len(content) > 200 {
			content = post.Content[:200]
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

		isFavorite, err := s.Dao.IsUserFavoritePost(req.UserId, post.Id)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrDatabase))
		}

		likeNum, err := s.Dao.GetLikedNum(item)
		if err != nil {
			logger.Error(err.Error(), zap.Error(errno.ErrRedis))
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
			Content:       content,
			CommentNum:    commentNum,
			LikeNum:       uint32(likeNum),
			IsLiked:       isLiked,
			IsFavorite:    isFavorite,
			Tags:          tags,
		}
	}

	return nil
}
