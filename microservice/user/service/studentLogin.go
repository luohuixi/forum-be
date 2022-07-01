package service

import (
	"context"
	"forum-user/dao"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/token"
)

// StudentLogin ... 登录
// 如果无 code，则返回 oauth 的地址，让前端去请求 oauth，
// 否则，用 code 获取 oauth 的 access token，并生成该应用的 auth token，返回给前端。
func (s *UserService) StudentLogin(ctx context.Context, req *pb.StudentLoginRequest, res *pb.LoginResponse) error {
	logger.Info("UserService StudentLogin")

	// 根据 StudentId 在 DB 查询 user
	user, err := s.Dao.GetUserByStudentId(req.StudentId)

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if user == nil {
		if err := auth.GetUserInfoFormOne(req.StudentId, req.Password); err != nil {
			return errno.ServerErr(errno.ErrRegister, err.Error())
		}
		info := &dao.RegisterInfo{
			StudentId: req.StudentId,
			Password:  req.Password,
			Role:      constvar.Normal,
		}
		// 用户未注册，自动注册
		if err := s.Dao.RegisterUser(info); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		// 注册后重新查询
		user, err = s.Dao.GetUserByStudentId(req.StudentId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
	} else {
		if !user.CheckPassword(req.Password) {
			return errno.ServerErr(errno.ErrPasswordIncorrect, "密码错误！")
		}
	}

	// 生成 auth token
	token, err := token.GenerateToken(&token.TokenPayload{
		Id:      user.Id,
		Role:    user.Role,
		Expired: util.GetExpiredTime(),
	})
	if err != nil {
		return errno.ServerErr(errno.ErrAuthToken, err.Error())
	}

	res.Token = token
	return nil
}
