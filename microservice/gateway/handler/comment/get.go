package comment

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

// Get ... 获取评论
// @Summary 获取评论 api
// @Description
// @Tags comment
// @Accept application/json
// @Produce application/json
// @Param comment_id path int true "comment_id"
// @Success 200 {object} Comment
// @Router /comment/{comment_id} [get]
func (a *Api) Get(c *gin.Context) {
	log.Info("Comment Get function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	getReq := &pb.Request{
		UserId: userId,
		Id:     uint32(id),
	}

	resp, err := service.PostClient.GetComment(context.TODO(), getReq)
	if err != nil {
		SendError(c, err, resp, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, resp)
}
