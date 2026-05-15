package service

import (
	"context"
	"errors"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

// Note: 目前评分是强一致性
// 后续可以在引入任务队列后改为最终一致性
// todo 将 平均分 单独拎出来异步刷新，crud 的时候都需要异步一下

func (s *PostService) CreateSipScoreEntryCommentRating(ctx context.Context, req *pb.CreateSipScoreEntryCommentRatingRequest, resp *pb.Response) error {
	logger.Info("PostService CreateSipScoreEntryCommentRating")

	rating := req.GetRating()
	if rating < 1 || rating > 5 {
		return errno.ServerErr(errno.ErrBadRequest, "rating must be between 1 and 5")
	}

	fc := func(tx *gorm.DB) error {
		// 锁住 entry，串行化同一 entry 的评分创建，避免并发下 participant_count 被放大。
		if err := s.Dao.LockSipScoreEntryForUpdate(req.GetSipScoreId(), req.GetSipScoreEntryId(), tx); err != nil {
			return err
		}

		existedRating, err := s.Dao.GetSipScoreEntryCommentRatingByUserForUpdate(req.GetSipScoreId(), req.GetSipScoreEntryId(), req.GetUserId(), tx)
		if err != nil {
			return err
		}
		if existedRating != nil {
			return gorm.ErrDuplicatedKey
		}

		// 1. 尝试创建评分记录，判断用户是否已经评分过
		data := &dao.SipScoreEntryCommentRating{
			CreatorID:      req.GetUserId(),
			LastModifiedBy: req.GetUserId(),
			SipScoreID:     req.GetSipScoreId(),
			EntryID:        req.GetSipScoreEntryId(),
			Rating:         rating,
			Content:        req.GetComment(),
			LikeNum:        0,
			ImgURL:         req.GetImgUrl(),
		}

		_, err = s.Dao.CreateSipScoreEntryCommentRating(data, tx)
		if err != nil {
			return err
		}

		// 2. 尝试更新 SipScore participant_count，判断 SipScore 是否存在
		err = s.Dao.IncrSipScoreParticipantCount(req.GetSipScoreId(), 1, tx)
		if err != nil {
			return err
		}

		// 3. 尝试更新 entry 统计字段，判断 entry 是否存在
		return s.Dao.IncrSipScoreEntryScore(req.GetSipScoreId(), req.GetSipScoreEntryId(), rating, 1, tx)
	}

	err := s.Dao.Transaction(fc)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errno.ServerErr(errno.ErrItemNotFound, "sip score or entry not found")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return errno.ServerErr(errno.ErrBadRequest, "user has already rated this entry")
	default:
		return errno.ServerErr(errno.ErrDatabase, "database transaction error")
	}
}
