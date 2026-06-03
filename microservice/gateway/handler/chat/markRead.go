package chat

import (
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"
	"strconv"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MarkRead ... 标记与某用户的私信为已读
// @Summary 标记与某用户的私信为已读
// @Tags chat
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param id path int true "聊天对象 id"
// @Success 200 {object} Response
// @Router /chat/read/{id} [patch]
func MarkRead(c *gin.Context) {
	log.Info("Chat MarkRead function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	id := c.Param("id")
	otherUserId, err := strconv.Atoi(id)
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, "id not legal", GetLine())
		return
	}
	if otherUserId <= 0 {
		SendError(c, errno.ErrPathParam, nil, "id not legal", GetLine())
		return
	}

	req := &pb.ReadRequest{
		UserId:      userId,
		OtherUserId: uint32(otherUserId),
	}
	resp, err := client.ChatClient.MarkRead(c.Request.Context(), req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, resp, nil)
}
