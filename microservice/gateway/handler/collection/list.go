package collection

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... list收藏
// @Summary list收藏 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []collection.Collection
// @Router /collection/list [get]
func (a *Api) List(c *gin.Context) {
	log.Info("Collection List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	listReq := &pb.UserIdRequest{
		UserId: userId,
	}

	resp, err := service.PostClient.ListCollection(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, resp, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.Collections)
}
