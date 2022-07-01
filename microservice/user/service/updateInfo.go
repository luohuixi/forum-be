package service

import (
	"context"

	pb "forum-user/proto"
	logger "forum/log"
	e "forum/pkg/err"
	errno "forum/pkg/err"
)

// UpdateInfo ... 更新用户信息
func (s *UserService) UpdateInfo(ctx context.Context, req *pb.UpdateInfoRequest, res *pb.Response) error {
	logger.Info("UserService UpdateInfo")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return e.ServerErr(errno.ErrUserExisted, err.Error())
	}

	user.Name = req.Info.Name
	user.Avatar = req.Info.AvatarUrl

	if err := user.Save(); err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
