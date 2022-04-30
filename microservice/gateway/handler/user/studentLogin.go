package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/pkg/errno"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	e "forum/pkg/err"

	errors "github.com/micro/go-micro/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StudentLogin ... 学生登录
// @Summary login api
// @Description login the student-forum
// @Tags auth
// @Accept  application/json
// @Produce  application/json
// @Param object body studentLoginRequest true "login_request"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /auth/login/student [post]
func StudentLogin(c *gin.Context) {
	log.Info("student login function called.",
		zap.String("X-Request-Id", util.GetReqID(c)))

	var req studentLoginRequest
	if err := c.Bind(&req); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	// 构造请求给 login
	loginReq := &pb.StudentLoginRequest{
		StudentId: req.StudentId,
		Password:  req.Password,
	}

	// 发送请求
	loginResp, err := service.UserClient.StudentLogin(context.Background(), loginReq)

	if err != nil {
		parsedErr := errors.Parse(err.Error())
		detail, errr := e.ParseDetail(parsedErr.Detail)

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
	resp := studentLoginResponse{
		Token: loginResp.Token,
	}

	SendResponse(c, nil, resp)
}
