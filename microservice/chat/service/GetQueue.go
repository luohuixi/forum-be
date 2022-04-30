package service

import (
	"context"
	pb "forum-chat/proto"
)

// GetQueue 获取队列
func (s *ChatService) GetQueue(ctx context.Context, req *pb.GetQueueRequest, resp *pb.GetQueueResponse) error {
	sum := sum256(req.UserId)
	// logger.With("tokenID", req.TokenId).With("sum", sum).Debugln("返回用户队列")
	// 在redis里面创建队列
	err := s.Dao.CreateQueue(sum)
	if err != nil {
		// logger.Errorln(err)
		return err
	}
	// resp.List = sum
	return err
}
