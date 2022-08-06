package chat

import (
	. "forum-gateway/handler"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"time"

	m "forum/model"
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
	userId := c.MustGet("userId")

	u4 := uuid.NewV4().String()

	if err := m.SetStringInRedis("user:"+u4, userId, time.Hour*24); err != nil {
		SendError(c, errno.ErrRedis, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, &Id{Id: u4})
}
