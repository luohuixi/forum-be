package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/pagetoken"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PostService) ListSipScoreEntryCommentRating(_ context.Context, req *pb.ListSipScoreEntryCommentRatingInfoRequest, resp *pb.ListSipScoreEntryCommentRatingInfoResponse) error {
	logger.Info("PostService ListSipScoreEntryCommentRating")

	if req.GetSipScoreId() == 0 || req.GetEntryId() == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id and entry_id required")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	var pageToken pb.SipScoreEntryCommentPageToken
	err := pagetoken.DecodePageToken(req.GetPageToken(), &pageToken)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
	}

	var pageTokenPtr *pb.SipScoreEntryCommentPageToken
	if req.GetPageToken() != "" {
		pageTokenPtr = &pageToken

		if pageTokenPtr.GetSortType() != req.GetSortType() {
			return errno.ServerErr(errno.ErrBadRequest, "page token sort type mismatch")
		}
	}

	switch req.GetSortType() {
	case constvar.SortByNewest:
		return s.listSipScoreEntryCommentRatingNewest(req.GetSipScoreId(), req.GetEntryId(), pageTokenPtr, limit, resp)
	case constvar.SortByHottest:
		return s.listSipScoreEntryCommentRatingHottest(req.GetSipScoreId(), req.GetEntryId(), pageTokenPtr, limit, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "sort type not legal")
	}
}

func (s *PostService) listSipScoreEntryCommentRatingCommon(
	sipScoreID, entryID uint32, pageToken *pb.SipScoreEntryCommentPageToken, limit uint32, resp *pb.ListSipScoreEntryCommentRatingInfoResponse,
	fetch func(token *pb.SipScoreEntryCommentPageToken, limit uint32) ([]*dao.SipScoreEntryCommentRating, error),
	nextToken func(last *dao.SipScoreEntryCommentRating) *pb.SipScoreEntryCommentPageToken,
) error {

	ratings, err := fetch(pageToken, limit)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.HasMore = false
	if uint32(len(ratings)) > limit-1 {
		resp.HasMore = true
		ratings = ratings[:limit-1]
	}

	if len(ratings) == 0 {
		resp.PageToken = ""
		resp.Ratings = nil
		return nil
	}

	if resp.HasMore {
		tokenStr, err := pagetoken.EncodePageToken(nextToken(ratings[len(ratings)-1]))
		if err != nil {
			return errno.ServerErr(errno.InternalServerError, "failed to encode page token")
		}
		resp.PageToken = tokenStr
	} else {
		resp.PageToken = ""
	}

	// 批量获取每个 rating 的前 3 条评论
	ratingIDs := make([]uint32, len(ratings))
	for i, r := range ratings {
		ratingIDs[i] = r.ID
	}

	commentsMap, err := s.Dao.BatchListCommentsByTargets(ratingIDs, constvar.SipScoreEntryCommentRating, 3)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Ratings = make([]*pb.SipScoreEntryCommentRating, len(ratings))
	for i, r := range ratings {
		resp.Ratings[i] = sipScoreEntryCommentRatingToPB(r, r.CommentNum, commentsMap[r.ID])
	}
	return nil
}

func (s *PostService) listSipScoreEntryCommentRatingNewest(sipScoreID, entryID uint32, pageToken *pb.SipScoreEntryCommentPageToken, limit uint32, resp *pb.ListSipScoreEntryCommentRatingInfoResponse) error {
	return s.listSipScoreEntryCommentRatingCommon(sipScoreID, entryID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryCommentPageToken, limit uint32) ([]*dao.SipScoreEntryCommentRating, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntryCommentRatingsNewest(sipScoreID, entryID, limit)
			}
			return s.Dao.ListSipScoreEntryCommentRatingsNewestWithCursor(
				sipScoreID, entryID,
				token.GetRatingId(),
				token.GetUpdatedAt().AsTime(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryCommentRating) *pb.SipScoreEntryCommentPageToken {
			return &pb.SipScoreEntryCommentPageToken{
				RatingId:  last.ID,
				UpdatedAt: timestamppb.New(last.UpdatedAt),
				SortType:  constvar.SortByNewest,
			}
		},
	)
}

func (s *PostService) listSipScoreEntryCommentRatingHottest(sipScoreID, entryID uint32, pageToken *pb.SipScoreEntryCommentPageToken, limit uint32, resp *pb.ListSipScoreEntryCommentRatingInfoResponse) error {
	return s.listSipScoreEntryCommentRatingCommon(sipScoreID, entryID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryCommentPageToken, limit uint32) ([]*dao.SipScoreEntryCommentRating, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntryCommentRatingsHottest(sipScoreID, entryID, limit)
			}
			return s.Dao.ListSipScoreEntryCommentRatingsHottestWithCursor(
				sipScoreID, entryID,
				token.GetRatingId(),
				token.GetLikeNum(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryCommentRating) *pb.SipScoreEntryCommentPageToken {
			return &pb.SipScoreEntryCommentPageToken{
				RatingId: last.ID,
				LikeNum:  last.LikeNum,
				SortType: constvar.SortByHottest,
			}
		},
	)
}

func sipScoreEntryCommentRatingToPB(r *dao.SipScoreEntryCommentRating, commentNum uint32, comments []*dao.CommentInfo) *pb.SipScoreEntryCommentRating {
	pbComments := make([]*pb.CommentInfo, len(comments))
	for i, c := range comments {
		pbComments[i] = &pb.CommentInfo{
			Id:            c.Id,
			TypeName:      c.TypeName,
			Content:       c.Content,
			FatherId:      c.FatherId,
			CreateTime:    timestamppb.New(c.CreateTime),
			CreatorId:     c.CreatorId,
			TargetId:      c.TargetID,
			TargetType:    c.TargetType,
			CreatorName:   c.CreatorName,
			CreatorAvatar: c.CreatorAvatar,
			LikeNum:       c.LikeNum,
			ImgUrl:        c.ImgUrl,
		}
	}

	return &pb.SipScoreEntryCommentRating{
		Id:              r.ID,
		SipScoreId:      r.SipScoreID,
		SipScoreEntryId: r.EntryID,
		CreatorId:       r.CreatorID,
		LastModifiedBy:  r.LastModifiedBy,
		Rating:          r.Rating,
		Content:         r.Content,
		LikeNum:         r.LikeNum,
		ImgUrl:          r.ImgURL,
		CreatedAt:       timestamppb.New(r.CreatedAt),
		UpdatedAt:       timestamppb.New(r.UpdatedAt),
		CommentNum:      commentNum,
		Comments:        pbComments,
	}
}
