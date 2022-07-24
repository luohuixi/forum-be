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

// Create ... 点赞
// @Summary 点赞 api
// @Description TypeId: Post = 1; Comment = 2
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body Item true "create_like_request"
// @Success 200 {object} handler.Response
// @Router /like [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Like Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req *pb.LikeItem
	if err := c.BindJSON(req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	likeReq := &pb.LikeRequest{
		UserId: userId,
		Item:   req,
	}

	_, err := service.PostClient.CreateLike(context.TODO(), likeReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}
