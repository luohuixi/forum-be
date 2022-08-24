package comment

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

// Create ... 创建评论
// @Summary 创建评论 api
// @Description typeName: first-level -> 一级评论; second-level -> 其它级
// @Tags comment
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_comment_request"
// @Success 200 {object} handler.Response
// @Router /comment [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Comment Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req pb.CreateCommentRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeName != constvar.FirstLevelComment && req.TypeName != constvar.SecondLevelComment {
		SendError(c, errno.ErrBadRequest, nil, "req.TypeName != constvar.FirstLevelComment && req.TypeName != constvar.SecondLevelComment", GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)
	req.CreatorId = userId

	ok, err := model.Enforce(req.CreatorId, constvar.Post, req.PostId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	resp, err := service.PostClient.CreateComment(context.TODO(), &req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 向 feed 发送请求
	pushReq := &pbf.PushRequest{
		Action: "创建",
		UserId: userId,
		Source: &pbf.Source{
			Id:       resp.Id,
			TypeName: constvar.Comment,
			Name:     req.Content,
		},
		TargetUserId: resp.TargetUserId,
	}
	_, err = service.FeedClient.Push(context.Background(), pushReq)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
