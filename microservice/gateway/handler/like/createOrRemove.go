package like

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
)

// CreateOrRemove ... 点赞/取消点赞
// @Summary 点赞/取消点赞 api
// @Description TypeName: post or comment
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body Item true "like_request"
// @Success 200 {object} handler.Response
// @Router /like [post]
func (a *Api) CreateOrRemove(c *gin.Context) {
	log.Info("Like CreateOrRemove function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req pb.LikeItem
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, req.TypeName, req.TargetId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	likeReq := &pb.LikeRequest{
		UserId: userId,
		Item:   &req,
	}

	_, err = service.PostClient.CreateOrRemoveLike(context.TODO(), likeReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
