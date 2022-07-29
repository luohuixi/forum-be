package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// GetInfo ... 获取用户信息
func (s *UserService) GetInfo(_ context.Context, req *pb.GetInfoRequest, resp *pb.UserInfoResponse) error {
	logger.Info("UserService GetInfo")

	list, err := s.Dao.GetUserByIds(req.Ids)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	userInfos := make([]*pb.UserInfo, len(list))
	for i, user := range list {
		userInfos[i] = &pb.UserInfo{
			Id:        user.Id,
			Name:      user.Name,
			AvatarUrl: user.Avatar,
			Email:     user.Email,
		}
	}

	resp.List = userInfos

	return nil
}
