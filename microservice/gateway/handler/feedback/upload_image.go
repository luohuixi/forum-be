package feedback

import (
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UploadImage ... 上传反馈截图
// @Summary 上传反馈截图 api
// @Tags feedback
// @Accept multipart/form-data
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param file formData file true "feedback_image"
// @Success 200 {object} Response{data=UploadImageResponse}
// @Router /feedback/image [post]
func (a *Api) UploadImage(c *gin.Context) {
	log.Info("Feedback UploadImage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	fileHeader, err := c.FormFile("file")
	if err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	fileToken, err := service.UploadFeedbackImage(c.Request.Context(), fileHeader)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, UploadImageResponse{FileToken: fileToken})
}
