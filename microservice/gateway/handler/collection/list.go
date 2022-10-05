package collection

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// List ... list收藏
// @Summary list收藏 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param user_id path int true "user_id"
// @Success 200 {object} []collection.Collection
// @Router /collection/list/{user_id} [get]
func (a *Api) List(c *gin.Context) {
	log.Info("Collection List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	targetUserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	if int(userId) != targetUserId {

		ok, err := model.Enforce(userId, constvar.CollectionAndLike, targetUserId, constvar.Read)
		if err != nil {
			SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
			return
		}

		if !ok {
			SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
			return
		}
	}

	listReq := &pb.UserIdRequest{
		UserId: uint32(targetUserId),
	}

	resp, err := service.PostClient.ListCollection(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, resp, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.Collections)
}
