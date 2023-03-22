package service

import (
	"context"
	"forum-chat/dao"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// Create 发送消息
func (s *ChatService) Create(_ context.Context, req *pb.CreateRequest, _ *pb.Response) error {
	logger.Info("CharService Create")

	data := &dao.ChatData{
		Content:  req.Content,
		Time:     req.Time,
		Receiver: req.TargetUserId,
		Sender:   req.UserId,
		TypeName: req.TypeName,
	}

	err := s.Dao.Create(data)

	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	return nil
}
