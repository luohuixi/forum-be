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
	"strconv"
)

// ListUserPost ... 获取用户发布的帖子
// @Summary list 用户发布的帖子 api
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param user_id path int true "user_id"
// @Success 200 {object} []post.Post
// @Router /post/published/{user_id} [get]
func (a *Api) ListUserPost(c *gin.Context) {
	log.Info("Post ListUserPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.Request{
		UserId: uint32(userId),
	}

	postResp, err := service.PostClient.ListUserCreatedPost(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, postResp.Posts)
}
