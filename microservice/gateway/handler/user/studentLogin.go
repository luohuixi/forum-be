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

// StudentLogin ... 学生登录
// @Summary login api
// @Description login the student-forum
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param object body studentLoginRequest true "login_request"
// @Success 200 {object} LoginResponse
// @Router /auth/login/student [post]
func StudentLogin(c *gin.Context) {
	log.Info("student login function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req studentLoginRequest
	if err := c.Bind(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
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
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 构造返回 response
	resp := studentLoginResponse{
		Token: loginResp.Token,
	}

	SendResponse(c, nil, resp)
}
