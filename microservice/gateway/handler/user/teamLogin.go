package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/pkg/errno"

	errors "github.com/micro/go-micro/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TeamLogin ... 团队登录
// @Summary login api
// @Description login the team-forum
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param object body teamLoginRequest true "login_request"
// @Success 200 {object} teamLoginResponse
// @Router /auth/login/team [post]
func TeamLogin(c *gin.Context) {
	log.Info("team login function called.",
		zap.String("X-Request-Id", util.GetReqID(c)))

	// 从前端获取 oauth_code
	var req teamLoginRequest
	if err := c.Bind(&req); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	// 构造请求给 login
	loginReq := &pb.TeamLoginRequest{
		OauthCode: req.OauthCode,
	}

	// 发送请求
	loginResp, err := service.UserClient.TeamLogin(context.Background(), loginReq)

	if err != nil {
		parsedErr := errors.Parse(err.Error())
		detail, errr := errno.ParseDetail(parsedErr.Detail)

		finalErrno := errno.InternalServerError
		if errr == nil {
			finalErrno = &errno.Errno{
				Code:    detail.Code,
				Message: detail.Cause,
			}
		}
		SendError(c, finalErrno, nil, err.Error(), GetLine())
		return
	}

	// 构造返回 response
	resp := teamLoginResponse{
		Token:       loginResp.Token,
		RedirectURL: loginResp.RedirectUrl,
	}

	SendResponse(c, nil, resp)
}
