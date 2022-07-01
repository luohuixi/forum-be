package service

import (
	"context"
	"forum/pkg/constvar"

	"forum-user/dao"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	logger "forum/log"
	e "forum/pkg/err"
	errno "forum/pkg/err"
	"forum/pkg/token"
)

// TeamLogin ... 登录
// 如果无 code，则返回 oauth 的地址，让前端去请求 oauth，
// 否则，用 code 获取 oauth 的 access token，并生成该应用的 auth token，返回给前端。
func (s *UserService) TeamLogin(ctx context.Context, req *pb.TeamLoginRequest, res *pb.LoginResponse) error {
	logger.Info("UserService TeamLogin")

	if req.OauthCode == "" {
		res.RedirectUrl = auth.OauthURL
		return nil
	}

	// get access token by auth code from auth-server
	if err := auth.OauthManager.ExchangeAccessTokenWithCode(req.OauthCode); err != nil {
		return e.ServerErr(errno.ErrRemoteAccessToken, err.Error())
	}

	// 尝试获取 access token，
	// 并在其中检查是否有效，如失效则尝试从 auth-server 更新
	accessToken, err := auth.OauthManager.GetAccessToken()
	if err != nil {
		return e.ServerErr(errno.ErrLocalAccessToken, err.Error())
	}

	// get user info by access token from auth-server
	userInfo, err := auth.GetInfoRequest(accessToken)
	if err != nil {
		return e.ServerErr(errno.ErrGetUserInfo, err.Error())
	}

	// 根据 email 在 DB 查询 user
	user, err := s.Dao.GetUserByEmail(userInfo.Email)

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	} else if user == nil {
		info := &dao.RegisterInfo{
			Name:  userInfo.Username,
			Email: userInfo.Email,
			Role:  constvar.TeamNormal,
		}
		// 用户未注册，自动注册
		if err := s.Dao.RegisterUser(info); err != nil {
			return e.ServerErr(errno.ErrDatabase, err.Error())
		}
		// 注册后重新查询
		user, err = s.Dao.GetUserByEmail(userInfo.Email)
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
