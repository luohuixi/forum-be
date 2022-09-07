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

// ListMainPost ... 获取主帖
// @Summary list 主帖 api
// @Description type_name : normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal); 根据category获取主帖list
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param type_name path string true "type_name"
// @Param category query string false "category"
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []post.Post
// @Router /post/list/{type_name} [get]
func (a *Api) ListMainPost(c *gin.Context) {
	log.Info("Post ListMainPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	typeName := c.Param("type_name")
	if typeName != constvar.NormalPost && typeName != constvar.MuxiPost {
		SendError(c, errno.ErrPathParam, nil, "type_name not legal", GetLine())
		return
	}

	ok, err := model.Enforce(userId, constvar.Post, typeName, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	category := c.DefaultQuery("category", "")

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

	listReq := &pb.ListMainPostRequest{
		UserId:     userId,
		Category:   category,
		TypeName:   typeName,
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
		Pagination: page != 0,
	}

	postResp, err := service.PostClient.ListMainPost(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, postResp.List)
}
