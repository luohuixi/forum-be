package service

import (
	"context"

	"forum-user/dao"
	errno "forum-user/errno"
	pb "forum-user/proto"
	e "forum/pkg/err"
)

// List ... 获取用户列表
func (s *UserService) List(ctx context.Context, req *pb.ListRequest, res *pb.ListResponse) error {

	// 过滤条件
	filter := &dao.UserModel{Role: req.Role}

	// DB 查询
	list, err := s.Dao.ListUser(req.Offset, req.Limit, req.LastId, filter)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	resList := make([]*pb.User, 0)

	for _, item := range list {
		resList = append(resList, &pb.User{
			Id:     item.Id,
			Name:   item.Name,
			Avatar: item.Avatar,
			Role:   item.Role,
		})
	}

	res.Count = uint32(len(list))
	res.List = resList

	return nil
}
