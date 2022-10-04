package service

import (
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"go.uber.org/zap"
)

// PostService ... 帖子服务
type PostService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *PostService {
	service := new(PostService)
	service.Dao = i
	return service
}

func (s PostService) processComments(userId uint32, comments []*pb.CommentInfo) {
	for _, comment := range comments {
		item := dao.Item{
			Id:       comment.Id,
			TypeName: constvar.Comment,
		}

		num, err := s.Dao.GetLikedNum(item)
		if err != nil {
			logger.Error(errno.ErrRedis.Error(), zap.String("cause", err.Error()))
		}
		comment.LikeNum = uint32(num)

		isLiked, err := s.Dao.IsUserHadLike(userId, item)
		if err != nil {
			logger.Error(errno.ErrRedis.Error(), zap.String("cause", err.Error()))
		}
		comment.IsLiked = isLiked
	}
}

func (s PostService) getPostInfo(postId uint32, userId uint32) (bool, bool, uint32, []string, uint32, uint32) {
	item := dao.Item{
		Id:       postId,
		TypeName: constvar.Post,
	}

	isLiked, err := s.Dao.IsUserHadLike(userId, item)
	if err != nil {
		logger.Error(errno.ErrRedis.Error(), zap.String("cause", err.Error()))
	}

	isCollection, err := s.Dao.IsUserCollectionPost(userId, postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), zap.String("cause", err.Error()))
	}

	likeNum, err := s.Dao.GetLikedNum(item)
	if err != nil {
		logger.Error(errno.ErrRedis.Error(), zap.String("cause", err.Error()))
	}

	tags, err := s.Dao.ListTagsByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), zap.String("cause", err.Error()))
	}

	commentNum, err := s.Dao.GetCommentNumByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), zap.String("cause", err.Error()))
	}

	collectionNum, err := s.Dao.GetCommentNumByPostId(postId)
	if err != nil {
		logger.Error(errno.ErrDatabase.Error(), zap.String("cause", err.Error()))
	}

	return isLiked, isCollection, uint32(likeNum), tags, commentNum, collectionNum
}
