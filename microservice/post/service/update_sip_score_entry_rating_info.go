package service

import (
	"context"
	"errors"
	pb "forum-post/proto"
	logger "forum/log"

	"forum/pkg/errno"
	"strings"

	"gorm.io/gorm"
)

func (s *PostService) UpdateSipScoreEntryCommentRatingInfo(_ context.Context, req *pb.UpdateSipScoreEntryCommentRatingInfoRequest, resp *pb.Response) error {
	logger.Info("PostService UpdateSipScoreEntryCommentRatingInfo")

	if req.GetSipScoreId() == 0 || req.GetSipScoreEntryId() == 0 || req.GetRatingId() == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id, sip_score_entry_id and rating_id required")
	}

	lastModifiedBy := req.GetLastModifiedBy()
	if lastModifiedBy == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "last_modified_by required")
	}

	updateMask := req.GetUpdateMask()
	if updateMask == nil || len(updateMask.Paths) == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "update_mask required")
	}

	field := map[string]interface{}{
		"rating":  req.GetRating(),
		"content": req.GetContent(),
		"img_url": req.GetImgUrl(),
	}

	update := map[string]interface{}{
		"last_modified_by": lastModifiedBy,
	}
	updateRating := false
	for _, path := range updateMask.Paths {
		if val, ok := field[path]; ok {
			if strings.Compare(path, "rating") == 0 {
				updateRating = true
				if req.GetRating() < 1 || req.GetRating() > 5 {
					return errno.ServerErr(errno.ErrBadRequest, "rating must be between 1 and 5")
				}
			}
			update[path] = val
		} else {
			return errno.ServerErr(errno.ErrBadRequest, "invalid update_mask path: "+path)
		}
	}

	fc := func(tx *gorm.DB) error {
		// 与创建评分保持一致：先锁 entry，再锁 rating，避免交叉顺序导致死锁。
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

		if updateRating {
			delta := int(req.GetRating()) - int(rating.Rating)
			if delta != 0 {
				if err := s.Dao.UpdateSipScoreEntryScoreByRatingDelta(req.GetSipScoreId(), req.GetSipScoreEntryId(), delta, tx); err != nil {
					return err
				}
			}
		}

		return s.Dao.UpdateSipScoreEntryCommentRating(req.GetSipScoreId(), req.GetSipScoreEntryId(), req.GetRatingId(), update, tx)
	}

	err := s.Dao.Transaction(fc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "sip_score_entry_comment_rating not found")
		}
		return errno.ServerErr(errno.ErrDatabase, "failed to update sip_score_entry_comment_rating: "+err.Error())
	}

	return nil
}
