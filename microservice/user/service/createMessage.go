package service

import (
	"context"

	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// CreateMessage ... 创建用户消息
func (s *UserService) CreateMessage(_ context.Context, req *pb.CreateMessageRequest, _ *pb.Response) error {
	logger.Info("UserService CreateMessage")

	if err := s.Dao.CreateMessage(0, req.Message); err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}
