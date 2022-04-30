package service

import (
	"context"
	pb "forum-chat/proto"
)

// GetList 获取列表
func (s *ChatService) GetList(ctx context.Context, req *pb.GetListRequest, resp *pb.GetListResponse) error {
	// 在redis里面获取该用户消息
	list, err := s.Dao.GetList(req.UserId)
	if err != nil {
		// logger.Errorln(err)
		return err
	}

	resp.List = list

	return err
}
