package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListUserCreatedPost(_ context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListUserCreatedPost")

	postIds, err := s.Dao.ListUserCreatedPost(req.UserId)
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
