package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListPopularTag ... 获取热门tags
// @Summary list 热门tags api
// @Description 降序
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []string
// @Router /post/popular_tag [get]
func (a *Api) ListPopularTag(c *gin.Context) {
	log.Info("Post ListPopularTag function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	resp, err := service.PostClient.ListPopularTag(context.TODO(), &pb.NullRequest{})
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.Tags)
}
