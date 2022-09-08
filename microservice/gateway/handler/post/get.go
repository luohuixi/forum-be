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
	"strconv"
)

// Get ... 获取帖子
// @Summary 获取帖子 api
// @Description
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param post_id path int true "post_id"
// @Success 200 {object} GetPostResponse
// @Router /post/{post_id} [get]
func (a *Api) Get(c *gin.Context) {
	log.Info("Post Get function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, id, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	getReq := &pb.Request{
		UserId: userId,
		Id:     uint32(id),
	}

	getResp, err := service.PostClient.GetPost(context.TODO(), getReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	var subPost []*SubPost

	for _, comment := range getResp.Comments {
		if comment.TypeName == constvar.SubPost {
			subPost = append(subPost, &SubPost{
				info: info{
					Id:      comment.Id,
					Content: comment.Content,
					// CommentNum:    comment.,
					Time:          comment.CreateTime,
					CreatorId:     comment.CreatorId,
					CreatorName:   comment.CreatorName,
					CreatorAvatar: comment.CreatorAvatar,
					LikeNum:       comment.LikeNum,
					IsLiked:       comment.IsLiked,
				},
				Comments: nil,
			})
		}
	}

	resp := GetPostResponse{
		info: info{
			Id:            getResp.Id,
			Content:       getResp.Content,
			CommentNum:    getResp.CommentNum,
			Time:          getResp.Time,
			CreatorId:     getResp.CreatorId,
			CreatorName:   getResp.CreatorName,
			CreatorAvatar: getResp.CreatorAvatar,
			LikeNum:       getResp.LikeNum,
			IsLiked:       getResp.IsLiked,
		},
		Title:        getResp.Title,
		Category:     getResp.Category,
		IsCollection: getResp.IsCollection,
		SubPosts:     subPost,
		Tags:         getResp.Tags,
	}

	// resp.setInfo(getResp)

	SendResponse(c, nil, resp)
}
