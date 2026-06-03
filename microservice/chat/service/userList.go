package service

import (
	"context"
	pb "forum-chat/proto"
)

// UserList ... 获取用户列表
func (s *ChatService) UserList(ctx context.Context, req *pb.UserListRequest, resp *pb.UserListResponse) error {
	userList, err := s.Dao.GetUserList(req.UserId, int(req.Limit), int(req.Page))
	if err != nil {
		return err
	}
	resp.UserLists = userList

	return nil
}
