package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *PostService) ListPopularTag(_ context.Context, _ *pb.NullRequest, resp *pb.Tags) error {
	logger.Info("PostService ListPopularTags")

	tags, err := s.Dao.ListPopularTags()
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp.Tags = tags

	return nil
}
