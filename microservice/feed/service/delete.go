package service

import (
	"context"
	pb "forum-feed/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

func (s *FeedService) Delete(_ context.Context, req *pb.Request, _ *pb.Response) error {
	logger.Info("FeedService Delete")

	err := s.Dao.Delete(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
