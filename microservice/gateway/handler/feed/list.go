package feed

import (
	"context"
	pb "forum-feed/proto"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... list feeds.
// @Summary list 此用户的动态 api
// @Tags feed
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Success 200 {object} FeedListResponse
// @Router /feed/list [get]
func (a *Api) List(c *gin.Context) {
	log.Info("feed List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)
	role := c.MustGet("role").(uint32)

	listReq := &pb.ListRequest{
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
		Role:       role,
		UserId:     userId,
		Pagination: limit != 0 || page != 0,
	}

	listResp, err := service.FeedClient.List(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, listResp)
}
