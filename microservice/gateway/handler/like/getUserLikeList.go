package like

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

// GetUserLikeList ... 获取用户点赞列表
// @Summary 获取用户点赞列表 api
// @Description TypeName: post or comment
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} ListResponse
// @Router /like/list [get]
func (a *Api) GetUserLikeList(c *gin.Context) { // TODO
	log.Info("Like GetUserLikeList function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	userIdReq := &pb.UserIdRequest{
		UserId: userId,
	}

	resp, err := service.PostClient.ListLikeByUserId(context.TODO(), userIdReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.List)
}
