package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	pbu "forum-user/proto"
	"forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... 获取帖子
// @Summary list post api
// @Description 获取帖子 (type_id = 1 -> 团队内(type_id暂时均填0))
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param type_id path int true "type_id"
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} post.ListResponse
// @Router /post/list/{type_id} [get]
func (a *Api) List(c *gin.Context) {
	log.Info("Post List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

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

	typeId := c.Param("type_id")

	userId := c.MustGet("userId").(uint32)

	if typeId != "" {
		ok, err := a.Dao.Enforce(userId, typeId, constvar.Read)
		if err != nil {
			SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
			return
		}

		if !ok {
			SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
			return
		}
	} else {
		typeId = "0"
	}

	listReq := &pb.ListPostRequest{
		LastId: uint32(lastId),
		Offset: uint32(page * limit),
		Limit:  uint32(limit),
		TypeId: typeId,
	}

	// 发送请求
	postResp, err := service.PostClient.ListPost(context.Background(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
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
		SendError(c, err, nil, "", GetLine())
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
