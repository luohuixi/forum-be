package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// GetProfile ... 获取用户个人信息
func (s *UserService) GetProfile(_ context.Context, req *pb.GetRequest, res *pb.UserProfile) error {
	logger.Info("UserService GetProfile")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, "")
	}

	res.Id = user.Id
	res.Name = user.Name
	res.Avatar = user.Avatar
	res.Email = user.Email
	res.Role = user.Role

	return nil
}
