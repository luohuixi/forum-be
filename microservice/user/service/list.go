package service

import (
	"context"

	"forum-user/dao"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// List ... 获取用户列表
func (s *UserService) List(_ context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {
	logger.Info("UserService List")

	// 过滤条件
	filter := &dao.UserModel{Role: req.Role}

	// DB 查询
	list, err := s.Dao.ListUser(req.Offset, req.Limit, req.LastId, filter)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resList := make([]*pb.User, len(list))

	for i, item := range list {
		resList[i] = &pb.User{
			Id:     item.Id,
			Name:   item.Name,
			Avatar: item.Avatar,
			Role:   item.Role,
		}
	}

	resp.Count = uint32(len(list))
	resp.List = resList

	return nil
}
