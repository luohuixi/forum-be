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

// ListSubPost ... 获取从帖
// @Summary list 从贴 api
// @Description type_id = 1 -> 团队内 (type_id暂时均填0); 根据 main_post_id 获取主贴的从贴list
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param main_post_id path int true "main_post_id"
// @Param object body ListSubPostRequest true "list_sub_post_request"
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} ListResponse
// @Router /post/list/{main_post_id} [get]
func (a *Api) ListSubPost(c *gin.Context) {
	log.Info("Post ListMainPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req *ListSubPostRequest
	if err := c.BindJSON(req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
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

	userId := c.MustGet("userId").(uint32)

	if req.TypeId != 0 {
		ok, err := a.Dao.Enforce(userId, req.TypeId, constvar.Read)
		if err != nil {
			SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
			return
		}

		if !ok {
			SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
			return
		}
	}

	listReq := &pb.ListSubPostRequest{
		UserId:     userId,
		MainPostId: req.MainPostId,
		TypeId:     req.TypeId,
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
	}

	if page != 0 {
		listReq.Pagination = true
	}

	postResp, err := service.PostClient.ListSubPost(c, listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	var ids []uint32
	for _, post := range postResp.List {
		ids = append(ids, post.CreatorId)
	}

	getReq := &pbu.GetInfoRequest{
		Ids: ids,
	}
	resp, err := service.UserClient.GetInfo(context.TODO(), getReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// TODO
	posts := make([]Post, len(resp.List))
	for i, user := range resp.List {
		posts[i] = Post{
			Content:       postResp.List[i].Content,
			Title:         postResp.List[i].Title,
			LastEditTime:  postResp.List[i].Time,
			Category:      postResp.List[i].Category,
			CreatorId:     postResp.List[i].CreatorId,
			IsLiked:       postResp.List[i].IsLiked,
			CreatorName:   user.Name,
			CreatorAvatar: user.AvatarUrl,
		}
	}

	SendResponse(c, errno.OK, ListResponse{posts: &posts})
}
