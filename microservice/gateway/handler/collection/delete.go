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

// Delete ... 取消收藏
// @Summary 取消收藏 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param collection_id path int true "collection_id"
// @Success 200 {object} handler.Response
// @Router /collection/{collection_id} [delete]
func (a *Api) Delete(c *gin.Context) {
	log.Info("Collection Delete function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	id, err := strconv.Atoi(c.Param("collection_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	ok, err := model.Enforce(userId, constvar.Collection, id, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	deleteReq := &pb.Request{
		Id:     uint32(id),
		UserId: userId,
	}

	_, err = service.PostClient.DeleteCollection(context.TODO(), deleteReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
