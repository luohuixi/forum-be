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
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListSubPost ... 获取从帖
// @Summary list 从帖 api
// @Description type_name : normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal); 根据 main_post_id 获取主帖的从帖list
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param type_name path string true "type_name"
// @Param main_post_id path int true "main_post_id"
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []post.Post
// @Router /post/list/{type_name}/{main_post_id} [get]
func (a *Api) ListSubPost(c *gin.Context) {
	log.Info("Post ListSubPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	typeName := c.Param("type_name")
	if typeName != constvar.NormalPost && typeName != constvar.MuxiPost {
		SendError(c, errno.ErrPathParam, nil, "type_name not legal", GetLine())
		return
	}

	mainPostId, err := strconv.Atoi(c.Param("main_post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	ok, err := model.Enforce(userId, constvar.Post, mainPostId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.ListSubPostRequest{
		UserId:     userId,
		MainPostId: uint32(mainPostId),
		TypeName:   typeName,
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
		Pagination: page != 0,
	}

	postResp, err := service.PostClient.ListSubPost(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, postResp.List)
}
