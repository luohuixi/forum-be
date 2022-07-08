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
)

// Create ... 创建评论
// @Summary 创建评论 api
// @Description
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest  true "create_comment_request"
// @Success 200 {object} handler.Response
// @Router /comment [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Comment Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req *pb.CreateCommentRequest
	if err := c.BindJSON(req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeId != 1 && req.TypeId != 2 {
		SendError(c, errno.ErrBadRequest, nil, "typeId != 1 && typeId != 2", GetLine())
		return
	}

	req.CreatorId = c.MustGet("userId").(uint32)

	// ok, err := a.Dao.Enforce(userId, typeId, constvar.Read)

	// 发送请求
	_, err := service.PostClient.CreateComment(context.Background(), req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}
