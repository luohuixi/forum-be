package agent

import (
	pb "forum-agent/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"

	"forum/client"
	"forum/pkg/errno"

	"forum/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GiveAnswer ... 生成帖子回复建议
// @Summary 生成帖子回复建议
// @Tags agent
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body GiveAnswerRequest true "give_answer_request"
// @Success 200 {object} handler.Response
// @Router /agent/answer [post]
func (a *Api) GiveAnswer(c *gin.Context) {
	log.Info("Agent GiveAnswer function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req GiveAnswerRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	giveReq := &pb.GiveAnswerRequest{PostId: req.PostId, ExtraContent: req.ExtraContent}
	if _, err := client.AgentClient.GiveAnswer(c.Request.Context(), giveReq); err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
