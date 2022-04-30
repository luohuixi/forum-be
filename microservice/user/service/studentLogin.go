package service

import (
	"context"
	"forum-user/errno"
	"forum-user/model"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	"forum/pkg/constvar"
	e "forum/pkg/err"
	"forum/pkg/token"
)

// StudentLogin ... 登录
// 如果无 code，则返回 oauth 的地址，让前端去请求 oauth，
// 否则，用 code 获取 oauth 的 access token，并生成该应用的 auth token，返回给前端。
func (s *UserService) StudentLogin(ctx context.Context, req *pb.StudentLoginRequest, res *pb.LoginResponse) error {
	// 根据 StudentId 在 DB 查询 user
	user, err := model.GetUserByStudentId(req.StudentId)

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	}
	if user == nil {
		if err := auth.GetUserInfoFormOne(req.StudentId, req.Password); err != nil {
			return e.ServerErr(errno.ErrRegister, err.Error())
		}
		info := &model.RegisterInfo{
			StudentId: req.StudentId,
			Password:  req.Password,
			Role:      constvar.Normal,
		}
		// 用户未注册，自动注册
		if err := model.RegisterUser(info); err != nil {
			return e.ServerErr(errno.ErrDatabase, err.Error())
		}
		// 注册后重新查询
		user, err = model.GetUserByStudentId(req.StudentId)
		if err != nil {
			return e.ServerErr(errno.ErrDatabase, err.Error())
		}
	}

	// 生成 auth token
	token, err := token.GenerateToken(&token.TokenPayload{
		Id:      user.Id,
		Role:    user.Role,
		Expired: util.GetExpiredTime(),
	})
	if err != nil {
		return e.ServerErr(errno.ErrAuthToken, err.Error())
	}

	res.Token = token
	return nil
}
