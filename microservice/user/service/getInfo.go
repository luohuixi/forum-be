package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// GetInfo ... 获取用户信息
func (s *UserService) GetInfo(_ context.Context, req *pb.GetInfoRequest, res *pb.UserInfoResponse) error {
	logger.Info("UserService GetInfo")

	list, err := s.Dao.GetUserByIds(req.Ids)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	var userInfos []*pb.UserInfo

	for _, user := range list {
		userInfos = append(userInfos, &pb.UserInfo{
			Id:        user.Id,
			Name:      user.Name,
			AvatarUrl: user.Avatar,
			Email:     user.Email,
		})
	}

	res.List = userInfos

	return nil
}
