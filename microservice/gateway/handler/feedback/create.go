package feedback

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Create ... 反馈与建议
// @Summary 反馈与建议 api
// @Tags feedback
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_feedback_request"
// @Success 200 {object} handler.Response
// @Router /feedback [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Feedback Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	_, err := client.PostClient.CreateFeedback(c.Request.Context(), &pb.CreateFeedbackRequest{
		UserId:   userID,
		Category: req.Category,
		Content:  req.Content,
		Contact:  req.Contact,
		ImgUrl:   req.ImgURL,
	})
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
