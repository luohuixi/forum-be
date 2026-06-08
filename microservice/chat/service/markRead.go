package service

import (
	"context"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// MarkRead 标记与指定用户的私信为已读
func (s *ChatService) MarkRead(_ context.Context, req *pb.ReadRequest, _ *pb.Response) error {
	logger.Info("CharService MarkRead")

	if err := s.Dao.MarkRead(req.UserId, req.OtherUserId); err != nil {
		return errno.ServerErr(errno.ErrListHistory, err.Error())
	}

	return nil
}
