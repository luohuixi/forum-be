package service

import (
	"context"
	"errors"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func (s *PostService) GetSipScoreEntryDetail(_ context.Context, req *pb.GetSipScoreEntryRequest, resp *pb.SipScoreEntryDetail) error {
	logger.Info("PostService GetSipScoreEntryDetail")

	if req.GetSipScoreId() == 0 || req.GetEntryId() == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id and entry_id required")
	}

	entry, err := s.Dao.GetSipScoreEntry(req.GetSipScoreId(), req.GetEntryId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NotFoundErr(errno.ErrItemNotFound, "sip-score-entry-"+strconv.Itoa(int(req.GetEntryId())))
		}
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Entry = sipScoreEntryModelToPB(entry)
	if req.GetUserId() != 0 {
		rating, err := s.Dao.GetSipScoreEntryCommentRatingByUser(req.GetSipScoreId(), req.GetEntryId(), req.GetUserId())
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if rating != nil {
			resp.MyRating = sipScoreEntryCommentRatingToPB(rating, rating.CommentNum, nil)
		}
	}

	return nil
}
