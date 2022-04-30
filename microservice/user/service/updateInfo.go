package service

import (
	"context"

	"forum-user/dao"
	"forum-user/errno"
	pb "forum-user/proto"
	e "forum/pkg/err"
)

// UpdateInfo ... 更新用户信息
func (s *UserService) UpdateInfo(ctx context.Context, req *pb.UpdateInfoRequest, res *pb.Response) error {
	user, err := dao.GetUser(req.Id)
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
