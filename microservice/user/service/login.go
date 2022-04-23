package service

import (
	"context"

	"forum-user/errno"
	"forum-user/model"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	e "forum/pkg/err"
	"forum/pkg/token"
)

// Login ... 登录
// 如果无 code，则返回 oauth 的地址，让前端去请求 oauth，
// 否则，用 code 获取 oauth 的 access token，并生成该应用的 auth token，返回给前端。
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest, res *pb.LoginResponse) error {
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

	// 根据 email 在本地 DB 查询 user
	user, err := model.GetUserByEmail(userInfo.Email)

	if err != nil {
		return e.ServerErr(errno.ErrDatabase, err.Error())
	} else if user == nil {
		info := &RegisterInfo{
			Name:  userInfo.Username,
			Email: userInfo.Email,
		}
		// 用户未注册，自动注册
		if err := RegisterUser(info); err != nil {
			return e.ServerErr(errno.ErrDatabase, err.Error())
		}
		// 注册后重新查询
		user, err = model.GetUserByEmail(userInfo.Email)
		if err != nil {
			return e.ServerErr(errno.ErrDatabase, err.Error())
		}
	}

	// 生成 auth token
	token, err := token.GenerateToken(&token.TokenPayload{
		ID:      user.ID,
		Role:    user.Role,
		TeamID:  user.TeamID,
		Expired: util.GetExpiredTime(),
	})
	if err != nil {
		return e.ServerErr(errno.ErrAuthToken, err.Error())
	}

	res.Token = token
	return nil
}
