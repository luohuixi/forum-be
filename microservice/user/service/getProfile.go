package service

import (
	"context"
	errno "forum-user/errno"
	"forum-user/model"
	pb "forum-user/proto"
	e "forum/pkg/err"
)

// GetProfile ... 获取用户个人信息
func (s *UserService) GetProfile(ctx context.Context, req *pb.GetRequest, res *pb.UserProfile) error {
	user, err := model.GetUser(req.Id)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return e.ServerErr(errno.ErrUserNotExisted, "")

	}
	res.Id = user.ID
	res.Name = user.Name
	res.Avatar = user.Avatar
	res.Email = user.Email
	res.Role = user.Role

	return nil
}
