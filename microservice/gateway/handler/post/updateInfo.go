package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/pkg/errno"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/pkg/constvar"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UpdateInfo ... 修改帖子信息
// @Summary update post info api
// @Description 修改帖子信息
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body UpdateInfoRequest  true "update_info_request"
// @Success 200 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /post [put]
func (a *Api) UpdateInfo(c *gin.Context) {
	log.Info("Post getInfo function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req UpdateInfoRequest
	if err := c.BindJSON(&req); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := a.Dao.Verify(userId, req.Id, constvar.Write)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendBadRequest(c, errno.ErrValidation, nil, "权限不足", GetLine())
		return
	}

	// 构造请求给 getInfo
	updateInfoReq := &pb.UpdateInfoRequest{
		Id:       req.Id,
		Content:  req.Content,
		Title:    req.Title,
		Category: req.Category,
		UserId:   userId,
	}

	// 发送请求
	_, err = service.PostClient.UpdateInfo(context.Background(), updateInfoReq)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}
