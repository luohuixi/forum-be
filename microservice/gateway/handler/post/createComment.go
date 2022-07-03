package post

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

// CreateComment ... 创建评论
// @Summary 创建评论 api
// @Description
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body ListRequest  true "list_request"
// @Success 200 {object} handler.ListResponse
// @Router /post/{type_id} [get]
func (a *Api) CreateComment(c *gin.Context) {
	log.Info("Create comment function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	var req *createCommentRequest
	if err := c.BindJSON(req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeId != 1 && req.TypeId != 2 {
		// ok, err := a.Dao.Enforce(userId, typeId, constvar.Read)

		SendError(c, errno.ErrBadRequest, nil, "typeId != 1 && typeId != 2", GetLine())
		return
	}

	// 构造请求给 list
	createReq := &pb.CreateCommentRequest{
		PostId:    req.PostId,
		TypeId:    req.TypeId,
		FatherId:  req.FatherId,
		Content:   req.Content,
		CreatorId: userId,
	}

	// 发送请求
	_, err := service.PostClient.CreateComment(context.Background(), createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}
