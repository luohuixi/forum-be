package chat

import (
	"context"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Id struct {
	Id string `json:"id"`
}

// GetId ... 获取该用户的uuid
// @Summary 获取该用户的uuid
// @Description 该用户发送信息前先获取自己的uuid，并放入query(id=?)，有效期24h
// @Tags chat
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} Id
// @Router /chat [get]
func GetId(c *gin.Context) {
	log.Info("Chat GetId function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	u4 := uuid.NewV4().String()

	req := pb.SetUUIdRequest{
		Uuid:   u4,
		UserId: userId,
	}

	_, err := service.ChatClient.SetUUId(context.TODO(), &req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, &Id{Id: u4})
}
