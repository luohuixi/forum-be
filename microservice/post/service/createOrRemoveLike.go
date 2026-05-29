package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) CreateOrRemoveLike(_ context.Context, req *pb.LikeRequest, _ *pb.Response) error {
	logger.Info("PostService CreateOrRemoveLike")

	switch req.Item.TypeName {
	case constvar.Post:
		return s.createOrRemovePostLike(req)
	case constvar.SipScoreEntryCommentRating:
		return s.createOrRemoveSipScoreEntryCommentRatingLike(req)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "target_type not legal")
	}
}

func (s *PostService) createOrRemovePostLike(req *pb.LikeRequest) error {
	var score int

	item := dao.Item{
		Id:       req.Item.TargetId,
		TypeName: req.Item.TypeName,
	}

	isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if isLiked {
		err = s.Dao.RemoveLike(req.UserId, item)
		score = -constvar.LikeScore
	} else {
		err = s.Dao.AddLike(req.UserId, item)
		score = constvar.LikeScore
	}
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	go func() {
		if err := s.Dao.AddChangeRecord(req.Item.TargetId); err != nil {
			logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
		}

		if err := s.Dao.ChangePostScore(req.Item.TargetId, score); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	return nil
}

// 创建/取消 sip-score-entry-comment-rating 点赞，创建时增加对应 rating 的点赞数，移除时减少
func (s *PostService) createOrRemoveSipScoreEntryCommentRatingLike(req *pb.LikeRequest) error {
	item := dao.Item{
		Id:       req.Item.TargetId,
		TypeName: req.Item.TypeName,
	}

	isLiked, err := s.Dao.IsUserHadLike(req.UserId, item)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if isLiked {
		if err = s.Dao.RemoveLike(req.UserId, item); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	} else {
		if err = s.Dao.AddLike(req.UserId, item); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	}

	rating, err := s.Dao.GetSipScoreEntryCommentRatingByID(req.Item.TargetId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if rating == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "rating not found")
	}

	if isLiked {
		if err := s.Dao.DecrSipScoreEntryCommentRatingLikeNum(rating.SipScoreID, rating.EntryID, rating.ID); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
	} else {
		if err := s.Dao.IncrSipScoreEntryCommentRatingLikeNum(rating.SipScoreID, rating.EntryID, rating.ID); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
	}

	return nil
}
