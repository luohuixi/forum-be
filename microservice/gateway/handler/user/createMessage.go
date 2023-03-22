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

// CreateMessage ... 创建 public message
// @Summary 创建 公共 message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateMessageRequest true "create_message_request"
// @Success 200 {object} handler.Response
// @Router /user/message [post]
func CreateMessage(c *gin.Context) {
	log.Info("User CreateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateMessageRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.CreateMessageRequest{
		Message: req.Message,
	}

	_, err := service.UserClient.CreateMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
