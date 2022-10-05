package collection

import (
	"context"
	pbf "forum-feed/proto"
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

// Create ... 收藏帖子
// @Summary 收藏帖子 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_collection_request"
// @Success 200 {object} handler.Response
// @Router /collection [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Collection Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, req.PostId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	createReq := pb.Request{
		UserId: userId,
		Id:     req.PostId,
	}

	resp, err := service.PostClient.CreateCollection(context.TODO(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 向 feed 发送请求
	pushReq := &pbf.PushRequest{
		Action: "收藏",
		UserId: userId,
		Source: &pbf.Source{
			Id:       resp.Id,
			TypeName: constvar.Collection,
			Name:     resp.TargetContent,
		},
		TargetUserId: resp.TargetUserId,
		Content:      "",
	}
	_, err = service.FeedClient.Push(context.TODO(), pushReq)

	SendResponse(c, err, nil)
}
