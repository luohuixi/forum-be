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

// Create ... 创建帖子
// @Summary 创建帖子 api
// @Description type_name = normal -> 团队外 (type_name暂时均填normal); 主贴的main_post_id不填或为0
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_post_request"
// @Success 200 {object} handler.Response
// @Router /post [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Post Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req *pb.CreatePostRequest
	if err := c.BindJSON(req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.MainPostId != 0 {
		req.CategoryId = 0
	}

	req.UserId = c.MustGet("userId").(uint32)

	_, err := service.PostClient.CreatePost(context.TODO(), req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}
