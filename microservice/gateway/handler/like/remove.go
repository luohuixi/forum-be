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
)

// Remove ... 取消点赞
// @Summary 取消点赞 api
// @Description TypeName: post or comment
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body Item true "remove_like_request"
// @Success 200 {object} handler.Response
// @Router /like [delete]
func (a *Api) Remove(c *gin.Context) {
	log.Info("Like Remove function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req pb.LikeItem
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	likeReq := &pb.LikeRequest{
		UserId: userId,
		Item:   &req,
	}

	_, err := service.PostClient.RemoveLike(context.TODO(), likeReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
