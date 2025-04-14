package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListMessage ... 获取 user message list
// @Summary list user message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} ListMessageResponse
// @Router /user/message/list [get]
func ListMessage(c *gin.Context) {
	log.Info("User ListMessage function called.", zap.String("X-Request-PostId", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	listReq := &pb.ListMessageRequest{
		UserId: userId,
	}

	listResp, err := service.UserClient.ListMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, ListMessageResponse{})
}
