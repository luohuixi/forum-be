package service

import (
	"context"

	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// ListMessage ... 获取用户消息列表
func (s *UserService) ListMessage(_ context.Context, req *pb.ListMessageRequest, resp *pb.ListMessageResponse) error {
	logger.Info("UserService ListMessage")

	// DB 查询
	messages, err := s.Dao.ListMessage(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp.Messages = messages

	return nil
}
