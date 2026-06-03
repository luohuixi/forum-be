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

// ListHistory ... 获取该用户的聊天记录
// @Summary 获取该用户的聊天记录
// @Tags chat
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param id path int true "id"
// @Success 200 {object} []Message
// @Router /chat/history/{id} [get]
func ListHistory(c *gin.Context) {
	log.Info("Chat ListHistory function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	id := c.Param("id")
	otherUserId, err := strconv.Atoi(id)
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, "id not legal", GetLine())
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	// 打开聊天记录前，先把双方 Redis 队列中尚未落库的消息同步进历史表。
	for _, pendingUserId := range []uint32{userId, uint32(otherUserId)} {
		getListRequest := &pb.GetListRequest{
			UserId: pendingUserId,
			Wait:   false,
		}

		_, err = client.ChatClient.GetList(c.Request.Context(), getListRequest)
		if err != nil {
			SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
			return
		}
	}

	req := pb.ListHistoryRequest{
		UserId:      userId,
		Offset:      uint32(page * limit),
		Limit:       uint32(limit),
		Pagination:  limit != 0 || page != 0,
		OtherUserId: uint32(otherUserId),
	}

	resp, err := client.ChatClient.ListHistory(c.Request.Context(), &req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.Histories)
}
