package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/log"
	"forum-gateway/pkg/errno"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	pbu "forum-user/proto"
	"forum/pkg/constvar"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... 获取帖子
// @Summary list post api
// @Description 获取帖子 (type = 1 -> 团队内)
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param Authorization header string true "token 用户令牌"
// @Param object body ListRequest  true "list_request"
// @Success 200 {object} handler.ListResponse
// @Router /post/{type_id} [get]
func (a *Api) List(c *gin.Context) {
	log.Info("Post list function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendBadRequest(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendBadRequest(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendBadRequest(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	typeId := c.Param("typeId")

	userId := c.MustGet("userId").(uint32)

	if typeId != "" {
		ok, err := a.Dao.Enforce(userId, typeId, constvar.Read)
		if err != nil {
			SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
			return
		}

		if !ok {
			SendBadRequest(c, errno.ErrValidation, nil, "权限不足", GetLine())
			return
		}
	} else {
		typeId = "0"
	}

	// 构造请求给 list
	listReq := &pb.ListRequest{
		LastId: uint32(lastId),
		Offset: uint32(page * limit),
		Limit:  uint32(limit),
		TypeId: typeId,
	}

	// 发送请求
	postResp, err := service.PostClient.List(context.Background(), listReq)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	var ids []uint32
	for _, post := range postResp.List {
		ids = append(ids, post.CreatorId)
	}

	req := &pbu.GetInfoRequest{
		Ids: ids,
	}
	resp, err := service.UserClient.GetInfo(context.Background(), req)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	var posts []Post
	for i, user := range resp.List {
		posts = append(posts, Post{
			Content:       postResp.List[i].Content,
			Title:         postResp.List[i].Title,
			LastEditTime:  postResp.List[i].Time,
			Category:      postResp.List[i].Category,
			CreatorId:     postResp.List[i].CreatorId,
			CreatorName:   user.Name,
			CreatorAvatar: user.AvatarUrl,
		})
	}

	SendResponse(c, errno.OK, ListResponse{posts: &posts})
}
