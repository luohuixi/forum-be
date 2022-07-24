package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TeamLogin ... 团队登录
// @Summary login api
// @Description login the team-forum
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param object body TeamLoginRequest true "login_request"
// @Success 200 {object} TeamLoginResponse
// @Router /auth/login/team [post]
func TeamLogin(c *gin.Context) {
	log.Info("User TeamLogin function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	// 从前端获取 oauth_code
	var req TeamLoginRequest
	if err := c.Bind(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	// 构造请求给 login
	loginReq := &pb.TeamLoginRequest{
		OauthCode: req.OauthCode,
	}

	loginResp, err := service.UserClient.TeamLogin(context.TODO(), loginReq)

	if err != nil {
		// parsedErr := errors.Parse(err.Error())
		// detail, errr := errno.ParseDetail(parsedErr.Detail)
		//
		// finalErrno := errno.InternalServerError
		// if errr == nil {
		// 	finalErrno = &errno.Errno{
		// 		Code:    detail.Code,
		// 		Message: detail.Cause,
		// 	}
		// }
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 构造返回 response
	resp := TeamLoginResponse{
		Token:       loginResp.Token,
		RedirectURL: loginResp.RedirectUrl,
	}

	SendResponse(c, nil, resp)
}
