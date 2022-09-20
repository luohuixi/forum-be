package post

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

// Create ... 创建帖子
// @Summary 创建帖子 api
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_post_request"
// @Success 200 {object} handler.Response
// @Router /post [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Post Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeName != constvar.NormalPost && req.TypeName != constvar.MuxiPost {
		SendError(c, errno.ErrBadRequest, nil, "type_name must be "+constvar.NormalPost+" or "+constvar.MuxiPost, GetLine())
		return
	}

	if req.ContentType != "md" && req.ContentType != "rtf" {
		SendError(c, errno.ErrBadRequest, nil, "content_type must be md or rtf", GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, req.TypeName, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	createReq := pb.CreatePostRequest{
		UserId:          userId,
		Content:         req.Content,
		TypeName:        req.TypeName,
		Title:           req.Title,
		Category:        req.Category,
		ContentType:     req.ContentType,
		Tags:            req.Tags,
		CompiledContent: req.CompiledContent,
		Summary:         req.Summary,
	}

	_, err = service.PostClient.CreatePost(context.TODO(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
