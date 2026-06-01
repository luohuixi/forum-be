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

func (s *PostService) SearchSipScoreEntry(_ context.Context, req *pb.SearchSipScoreEntryRequest, resp *pb.ListSipScoreEntryResponse) error {
	logger.Info("PostService SearchSipScoreEntry")

	keyword := req.GetKeyword()
	if keyword == "" {
		return errno.ServerErr(errno.ErrBadRequest, "keyword required")
	}

	sipScoreID := req.GetSipScoreId()
	if sipScoreID == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id required")
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	var entries []*dao.SipScoreEntryModel
	var err error

	if req.GetPageToken() == "" {
		entries, err = s.Dao.SearchSipScoreEntry(sipScoreID, keyword, limit)
	} else {
		var pageToken pb.SipScoreEntryPageToken
		if err = pagetoken.DecodePageToken(req.GetPageToken(), &pageToken); err != nil {
			return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
		}

		entries, err = s.Dao.SearchSipScoreEntryWithCursor(
			sipScoreID,
			keyword,
			pageToken.GetEntryId(),
			pageToken.GetUpdatedAt().AsTime(),
			limit,
		)
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	hasMore := len(entries) == int(limit)
	if hasMore {
		entries = entries[:len(entries)-1]
	}

	resp.Entries = make([]*pb.SipScoreEntry, len(entries))
	for i, e := range entries {
		resp.Entries[i] = sipScoreEntryModelToPB(e)
	}

	if hasMore {
		last := entries[len(entries)-1]
		nextToken := &pb.SipScoreEntryPageToken{
			EntryId:   last.ID,
			UpdatedAt: timestamppb.New(last.UpdatedAt),
		}
		encoded, err := pagetoken.EncodePageToken(nextToken)
		if err != nil {
			logger.Error("encode search entry page token error", logger.String(err.Error()))
		} else {
			resp.PageToken = encoded
		}
	}
	resp.HasMore = hasMore

	return nil
}
