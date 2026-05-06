package agent

import (
	pb "forum-agent/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"time"

	"github.com/gin-gonic/gin"
	mclient "go-micro.dev/v4/client"
	"go.uber.org/zap"
)

// AddKnowledge ... 将帖子加入知识库
// @Summary 将帖子加入知识库
// @Tags agent
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body AddKnowledgeRequest true "add_knowledge_request"
// @Success 200 {object} handler.Response
// @Router /agent/knowledge [post]
func (a *Api) AddKnowledge(c *gin.Context) {
	log.Info("Agent AddKnowledge function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	ok, err := model.HasRole(userId, constvar.NormalAdminRole)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "permission denied", GetLine())
	}

	var req AddKnowledgeRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	addReq := &pb.AddKnowledgeRequest{PostId: req.PostId}
	_, err = client.AgentClient.AddKnowledge(c.Request.Context(),
		addReq,
		mclient.WithRequestTimeout(5*time.Minute),
		withConnectionTimeout(5*time.Minute),
	)

	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}

func withConnectionTimeout(d time.Duration) mclient.CallOption {
	return func(o *mclient.CallOptions) {
		o.ConnectionTimeout = d
	}
}
