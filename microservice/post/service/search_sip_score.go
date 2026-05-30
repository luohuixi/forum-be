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

func (s *PostService) SearchSipScore(ctx context.Context, req *pb.SearchSipScoreRequest, resp *pb.SearchSipScoreResponse) error {
	logger.Info("PostService SearchSipScore")

	keyword := req.GetKeyword()
	if keyword == "" {
		return errno.ServerErr(errno.ErrBadRequest, "keyword required")
	}

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	domain := s.resolveListDomain(ctx, req.GetUserId())

	var sipScores []*dao.SipScoreModel
	var err error

	if req.GetPageToken() == "" {
		sipScores, err = s.Dao.SearchSipScore(keyword, limit, domain)
	} else {
		var pageToken pb.SipScorePageToken
		if err = pagetoken.DecodePageToken(req.GetPageToken(), &pageToken); err != nil {
			return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
		}

		sipScores, err = s.Dao.SearchSipScoreWithCursor(
			keyword,
			pageToken.GetId(),
			pageToken.GetUpdatedAt().AsTime(),
			limit,
			domain,
		)
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	hasMore := len(sipScores) == int(limit)
	if hasMore {
		sipScores = sipScores[:len(sipScores)-1]
	}

	// 批量拉取每个榜单的最热 entries
	ids := make([]uint32, len(sipScores))
	for i, s := range sipScores {
		ids[i] = s.ID
	}
	entriesMap, _ := s.Dao.BatchListSipScoreEntriesHottest(ids, constvar.DefaultPageSize)

	if hasMore {
		last := sipScores[len(sipScores)-1]
		nextToken := &pb.SipScorePageToken{
			Id:        last.ID,
			UpdatedAt: timestamppb.New(last.UpdatedAt),
		}
		encoded, err := pagetoken.EncodePageToken(nextToken)
		if err != nil {
			logger.Error("encode search page token error", logger.String(err.Error()))
		} else {
			resp.PageToken = encoded
		}
	}
	resp.HasMore = hasMore

	resp.SipScores = make([]*pb.SipScoreWithEntries, len(sipScores))
	for i, m := range sipScores {
		resp.SipScores[i] = &pb.SipScoreWithEntries{
			Meta:    sipScoreModelToPB(m, false),
			Entries: sipScoreEntriesModelToPB(entriesMap[m.ID]),
		}
	}

	return nil
}
