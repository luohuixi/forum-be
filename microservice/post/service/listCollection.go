package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListCollection(_ context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListCollections")

	postIds, err := s.Dao.ListCollectionByUserId(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	posts, err := s.Dao.ListPostInfoByPostIds(postIds, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrListPostInfoByPostIds, err.Error())
	}

	resp.Posts = posts

	return nil
}
