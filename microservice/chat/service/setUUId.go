package service

import (
	"context"
	pb "forum-chat/proto"
	logger "forum/log"
	m "forum/model"
	"forum/pkg/errno"
	"time"
)

// SetUUId 存uuId到redis
func (s *ChatService) SetUUId(_ context.Context, req *pb.SetUUIdRequest, _ *pb.Response) error {
	logger.Info("CharService SetUUId")

	if err := m.SetStringInRedis("user:"+req.Uuid, req.UserId, time.Hour*24); err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}
