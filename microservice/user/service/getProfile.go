package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	e "forum/pkg/err"
	errno "forum/pkg/err"
)

// GetProfile ... 获取用户个人信息
func (s *UserService) GetProfile(ctx context.Context, req *pb.GetRequest, res *pb.UserProfile) error {
	logger.Info("UserService GetProfile")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return e.ServerErr(errno.ErrUserNotExisted, "")

	}
	res.Id = user.Id
	res.Name = user.Name
	res.Avatar = user.Avatar
	res.Email = user.Email
	res.Role = user.Role

	return nil
}
