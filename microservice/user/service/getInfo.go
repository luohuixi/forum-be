package service

import (
	"context"
	errno "forum-user/errno"
	pb "forum-user/proto"
	logger "forum/log"
	e "forum/pkg/err"
)

// GetInfo ... 获取用户信息
func (s *UserService) GetInfo(ctx context.Context, req *pb.GetInfoRequest, res *pb.UserInfoResponse) error {
	logger.Info("UserService GetInfo")

	list, err := s.Dao.GetUserByIds(req.Ids)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	userInfos := make([]*pb.UserInfo, 0)

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
