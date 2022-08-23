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

// UpdateInfo ... 修改帖子信息
// @Summary update post info api
// @Description 修改帖子信息
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body UpdateInfoRequest  true "update_info_request"
// @Success 200 {object} handler.Response
// @Router /post [put]
func (a *Api) UpdateInfo(c *gin.Context) {
	log.Info("Post UpdateInfo function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req pb.UpdatePostInfoRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	req.UserId = c.MustGet("userId").(uint32)

	ok, err := model.Enforce(req.UserId, constvar.Post, req.Id, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	_, err = service.PostClient.UpdatePostInfo(context.TODO(), &req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
