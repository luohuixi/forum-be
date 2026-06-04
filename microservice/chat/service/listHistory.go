package service

import (
	"context"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// ListHistory 获取聊天记录
func (s *ChatService) ListHistory(_ context.Context, req *pb.ListHistoryRequest, resp *pb.ListHistoryResponse) error {
	logger.Info("CharService ListHistories")

	// get message histories of the user from redis
	histories, err := s.Dao.ListHistory(req.UserId, req.OtherUserId, req.Offset, req.Limit, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrListHistory, err.Error())
	}

	resp.Histories = histories

	return nil
}
