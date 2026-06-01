package service

import (
	"context"
	"errors"

	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

func (s *PostService) DeleteSipScoreEntryCommentRating(_ context.Context, req *pb.DeleteSipScoreEntryCommentRatingRequest, _ *pb.Response) error {
	logger.Info("PostService DeleteSipScoreEntryCommentRating")

	if req.GetSipScoreId() == 0 || req.GetSipScoreEntryId() == 0 || req.GetRatingId() == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id, sip_score_entry_id and rating_id required")
	}

	fc := func(tx *gorm.DB) error {
		// 先锁 entry，串行化同一 entry 操作
		if err := s.Dao.LockSipScoreEntryForUpdate(req.GetSipScoreId(), req.GetSipScoreEntryId(), tx); err != nil {
			return err
		}

		rating, err := s.Dao.GetSipScoreEntryCommentRatingForUpdate(req.GetSipScoreId(), req.GetSipScoreEntryId(), req.GetRatingId(), tx)
		if err != nil {
			return err
		}
		if rating == nil {
			return gorm.ErrRecordNotFound
		}

		// 1. 删除评分记录
		if err := s.Dao.DeleteSipScoreEntryCommentRating(req.GetSipScoreId(), req.GetSipScoreEntryId(), req.GetRatingId(), tx); err != nil {
			return err
		}

		// 2. SipScore 参与人数 -1
		if err := s.Dao.IncrSipScoreParticipantCount(req.GetSipScoreId(), -1, tx); err != nil {
			return err
		}

		// 3. Entry 统计回退：score_total -= rating，participant_count -= 1，重算 score_avg
		return s.Dao.DecrSipScoreEntryScore(req.GetSipScoreId(), req.GetSipScoreEntryId(), rating.Rating, tx)
	}

	err := s.Dao.Transaction(fc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "rating not found")
		}
		return errno.ServerErr(errno.ErrDatabase, "failed to delete rating: "+err.Error())
	}

	// todo 待日后异步删除这条 rating 下面的回复内容
	return nil
}
