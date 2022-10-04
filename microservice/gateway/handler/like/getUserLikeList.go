package like

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// GetUserLikeList ... 获取用户点赞列表
// @Summary 获取用户点赞列表 api
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param user_id path int true "user_id"
// @Success 200 {object} ListResponse
// @Router /like/list/{user_id} [get]
func (a *Api) GetUserLikeList(c *gin.Context) { // TODO
	log.Info("Like GetUserLikeList function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}
	userIdReq := &pb.UserIdRequest{
		UserId: uint32(userId),
	}

	resp, err := service.PostClient.ListLikeByUserId(context.TODO(), userIdReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.List)
}
