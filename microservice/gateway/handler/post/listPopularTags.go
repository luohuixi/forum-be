package post

import (
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListPopularTags ... 获取热门tags
// @Summary list 热门tags api
// @Description 降序
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []string
// @Router /post/tags [get]
func (a *Api) ListPopularTags(c *gin.Context) {
	log.Info("Post ListPopularTags function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	resp, err := service.PostClient.ListPopularTags(c, &pb.NullRequest{})
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, resp.Tags)
}
