package sipscore

import (
	pb "forum-post/proto"
	"time"
)

func sipScoreEntryFromPB(rpcEntry *pb.SipScoreEntry) *SipScoreEntry {
	if rpcEntry == nil {
		rpcEntry = &pb.SipScoreEntry{}
	}
	return &SipScoreEntry{
		ID:               rpcEntry.Id,
		SipScoreID:       rpcEntry.SipScoreId,
		CreatedAt:        rpcEntry.CreatedAt.AsTime().Format(time.DateTime),
		UpdatedAt:        rpcEntry.UpdatedAt.AsTime().Format(time.DateTime),
		Creator:          &userInfo{ID: rpcEntry.CreatorId},
		LastModifiedBy:   &userInfo{ID: rpcEntry.LastModifiedBy},
		Name:             rpcEntry.Name,
		Description:      rpcEntry.Description,
		CoverImg:         rpcEntry.CoverImg,
		ParticipantCount: rpcEntry.ParticipantCount,
		CommentCount:     rpcEntry.CommentCount,
		ScoreTotal:       rpcEntry.ScoreTotal,
		ScoreAvg:         rpcEntry.ScoreAvg,
	}
}

func sipScoreEntriesFromPB(rpcEntries []*pb.SipScoreEntry) []*SipScoreEntry {
	httpEntries := make([]*SipScoreEntry, len(rpcEntries))
	for i, rpcEntry := range rpcEntries {
		httpEntries[i] = sipScoreEntryFromPB(rpcEntry)
	}
	return httpEntries
}

func sipScoreFromPB(meta *pb.SipScore) *SipScore {
	if meta == nil {
		meta = &pb.SipScore{}
	}
	return &SipScore{
		ID:               meta.Id,
		CreatedAt:        meta.CreatedAt.AsTime().Format(time.DateTime),
		UpdatedAt:        meta.UpdatedAt.AsTime().Format(time.DateTime),
		Creator:          &userInfo{ID: meta.CreatorId},
		LastModifiedBy:   &userInfo{ID: meta.LastModifiedBy},
		EntryCount:       meta.EntryCount,
		CollectCount:     meta.CollectCount,
		ParticipantCount: meta.ParticipantCount,
		Name:             meta.Name,
		Description:      meta.Description,
		CoverImg:         meta.CoverImg,
		Domain:           meta.Domain,
		Category:         meta.Category,
		Tags:             meta.Tags,
		IsCollected:      meta.IsCollected,
	}
}

func sipScoreWithEntriesFromPB(rpcSipScore *pb.SipScoreWithEntries) *SipScoreWithEntries {
	if rpcSipScore == nil {
		rpcSipScore = &pb.SipScoreWithEntries{}
	}
	return &SipScoreWithEntries{
		SipScore: sipScoreFromPB(rpcSipScore.GetMeta()),
		Entries:  sipScoreEntriesFromPB(rpcSipScore.GetEntries()),
	}
}

func listSipScoresResponseFromPB(resp *pb.ListSipScoreResponse) *ListSipScoresResponse {
	if resp == nil {
		resp = &pb.ListSipScoreResponse{}
	}
	items := make([]*SipScoreWithEntries, len(resp.GetSipScores()))
	for i, rpcSipScore := range resp.GetSipScores() {
		items[i] = sipScoreWithEntriesFromPB(rpcSipScore)
	}
	return &ListSipScoresResponse{
		SipScores: items,
		PageToken: resp.PageToken,
		HasMore:   resp.HasMore,
	}
}

func sipScoreEntryRatingFromPB(r *pb.SipScoreEntryCommentRating) *SipScoreEntryCommentRatingInfo {
	if r == nil || r.Id == 0 {
		return nil
	}
	return &SipScoreEntryCommentRatingInfo{
		ID:              r.Id,
		SipScoreID:      r.SipScoreId,
		SipScoreEntryID: r.SipScoreEntryId,
		Creator:         &userInfo{ID: r.CreatorId},
		LastModifiedBy:  &userInfo{ID: r.LastModifiedBy},
		Rating:          r.Rating,
		Content:         r.Content,
		CommentID:       r.CommentId,
		LikeNum:         r.LikeNum,
		ImgUrl:          r.ImgUrl,
		CreatedAt:       r.CreatedAt.AsTime().Format(time.DateTime),
		UpdatedAt:       r.UpdatedAt.AsTime().Format(time.DateTime),
		CommentNum:      r.CommentNum,
		Comments:        []*CommentInfo{},
	}
}
