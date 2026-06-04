package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// DeletePrivateMessage ... 删除用户消息
func (s *UserService) DeletePrivateMessage(_ context.Context, req *pb.DeletePrivateMessageRequest, _ *pb.Response) error {
	logger.Info("UserService DeletePrivateMessage")

	if req.Id == "" {
		if err := s.Dao.DeleteMessage(req.UserId); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	} else {
		if err := s.Dao.DeleteOneMessage(req.UserId, req.Id); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	}

	return nil
}

// MarkPrivateMessageRead ... 标记用户消息为已读
func (s *UserService) MarkPrivateMessageRead(_ context.Context, req *pb.DeletePrivateMessageRequest, _ *pb.Response) error {
	logger.Info("UserService MarkPrivateMessageRead")

	if req.Id == "" {
		if err := s.Dao.MarkAllMessageRead(req.UserId); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	} else {
		if err := s.Dao.MarkOneMessageRead(req.UserId, req.Id); err != nil {
			return errno.ServerErr(errno.ErrRedis, err.Error())
		}
	}

	return nil
}
