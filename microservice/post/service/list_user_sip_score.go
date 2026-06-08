package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListUserCreatedSipScore(ctx context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListSipScoreResponse) error {
	logger.Info("PostService ListUserCreatedSipScore")

	targetUserID := req.GetTargetUserId()
	if targetUserID == 0 {
		targetUserID = req.GetUserId()
	}
	sipScores, err := s.Dao.ListSipScoreByCreator(targetUserID, req.GetOffset(), req.GetLimit(), req.GetLastId(), req.GetPagination())
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	return s.fillSipScoreListResponse(ctx, req.GetUserId(), sipScores, resp)
}

func (s *PostService) ListUserCollectedSipScore(ctx context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListSipScoreResponse) error {
	logger.Info("PostService ListUserCollectedSipScore")

	targetUserID := req.GetTargetUserId()
	if targetUserID == 0 {
		targetUserID = req.GetUserId()
	}
	sipScores, err := s.Dao.ListCollectedSipScoreByUser(targetUserID, req.GetOffset(), req.GetLimit(), req.GetLastId(), req.GetPagination())
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	return s.fillSipScoreListResponse(ctx, req.GetUserId(), sipScores, resp)
}

func (s *PostService) fillSipScoreListResponse(_ context.Context, viewerID uint32, sipScores []*dao.SipScoreModel, resp *pb.ListSipScoreResponse) error {
	if len(sipScores) == 0 {
		resp.SipScores = []*pb.SipScoreWithEntries{}
		return nil
	}

	ids := make([]uint32, len(sipScores))
	for i, sip := range sipScores {
		ids[i] = sip.ID
	}

	collected, err := s.Dao.ListIsUserCollected(viewerID, constvar.CollectionSipScore, ids)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	entries, err := s.Dao.BatchListSipScoreEntriesHottest(ids, constvar.DefaultPageSize)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.SipScores = make([]*pb.SipScoreWithEntries, len(sipScores))
	for i, sip := range sipScores {
		resp.SipScores[i] = &pb.SipScoreWithEntries{
			Meta:    sipScoreModelToPB(sip, collected[sip.ID]),
			Entries: sipScoreEntriesModelToPB(entries[sip.ID]),
		}
	}
	return nil
}
