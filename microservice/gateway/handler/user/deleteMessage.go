package user

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func privateMessageRequest(c *gin.Context) *pb.DeletePrivateMessageRequest {
	userId := c.MustGet("userId").(uint32)
	return &pb.DeletePrivateMessageRequest{
		UserId: userId,
		Id:     c.Query("id"),
	}
}

func markPrivateMessageRead(c *gin.Context) {
	_, err := client.UserClient.MarkPrivateMessageRead(c.Request.Context(), privateMessageRequest(c))
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}

// ReadPrivateMessage ... 标记 private message 已读
// @Summary 标记 个人 message 已读 api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /user/private_message/read [patch]
func ReadPrivateMessage(c *gin.Context) {
	log.Info("User ReadPrivateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	markPrivateMessageRead(c)
}

// DeletePrivateMessage ... 删除 private message
// @Summary 删除 个人 message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /user/private_message [delete]
func DeletePrivateMessage(c *gin.Context) {
	log.Info("User DeletePrivateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	_, err := client.UserClient.DeletePrivateMessage(c.Request.Context(), privateMessageRequest(c))
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
