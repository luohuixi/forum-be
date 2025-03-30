package service

import (
	"context"
	"forum-user/dao"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/token"
)

// StudentLogin ... 登录
// 如果无 code，则返回 oauth 的地址，让前端去请求 oauth，
// 否则，用 code 获取 oauth 的 access token，并生成该应用的 auth token，返回给前端。
func (s *UserService) StudentLogin(_ context.Context, req *pb.StudentLoginRequest, resp *pb.LoginResponse) error {
	logger.Info("UserService StudentLogin")

	// 使用 ccnu 登陆
	if err := auth.GetUserInfoFormOne(req.StudentId, req.Password); err != nil {
		return errno.ServerErr(errno.ErrPasswordIncorrect, err.Error())
	}

	//查询是否存在用户
	user, err := s.Dao.GetUserByStudentId(req.StudentId)
	if err != nil {
		return err
	}

	//如果用户为空
	if user == nil {
		info := &dao.RegisterInfo{
			StudentId: req.StudentId,
			Password:  req.Password,
			Role:      constvar.NormalRole,
			Name:      req.StudentId,
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

		if err := s.Dao.AddPublicPolicy(user.Role, user.Id); err != nil {
			return errno.ServerErr(errno.ErrCasbin, err.Error())
		}

		if err := model.AddRole("user", user.Id, constvar.NormalRole); err != nil {
			return errno.ServerErr(errno.ErrCasbin, err.Error())
		}
	} else {
		//更新用户密码
		err := s.Dao.UpdatePassword(user.Id, req.Password)
		if err != nil {
			return err
		}
	}

	//根据权限生成token
	role := uint32(constvar.Normal)
	if user.Role == constvar.NormalAdminRole || user.Role == constvar.MuxiAdminRole {
		role = constvar.Admin
	} else if user.Role == constvar.SuperAdminRole {
		role = constvar.SuperAdmin
	}

	// 生成 auth token
	Token, err := token.GenerateToken(&token.TokenPayload{
		Id:      user.Id,
		Role:    role,
		Expired: util.GetExpiredTime(),
	})
	if err != nil {
		return errno.ServerErr(errno.ErrAuthToken, err.Error())
	}

	resp.Token = Token
	return nil
}
