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
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param post_id path int true "post_id"
// @Param Authorization header string true "token 用户令牌"
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

	// 三次遍历分别处理 SubPost FirstLevelComment SecondLevelComment
	var subPosts []*SubPost
	subPostCommentsMap := make(map[uint32]*[]*Comment)

	firstLevelCommentsMap := make(map[uint32]*[]*info)

	for _, comment := range getResp.Comments {
		if comment.TypeName == constvar.SubPost {
			var subPost = new(SubPost)

			subPostCommentsMap[comment.Id] = &subPost.Comments

			setInfo(&subPost.info, *comment)

			subPosts = append(subPosts, subPost)

		}
	}

	for _, comment := range getResp.Comments {
		if comment.TypeName == constvar.FirstLevelComment {
			var subPostComment = new(Comment)

			subPostComments := subPostCommentsMap[comment.FatherId]

			firstLevelCommentsMap[comment.Id] = &subPostComment.Replies

			setInfo(&subPostComment.info, *comment)

			*subPostComments = append(*subPostComments, subPostComment)
		}
	}

	for _, comment := range getResp.Comments {
		if comment.TypeName == constvar.SecondLevelComment {
			var commentReply = new(info)

			commentsReplies := firstLevelCommentsMap[comment.FatherId]

			setInfo(commentReply, *comment)

			*commentsReplies = append(*commentsReplies, commentReply)
		}
	}

	resp := GetPostResponse{
		Title:        getResp.Title,
		Category:     getResp.Category,
		IsCollection: getResp.IsCollection,
		SubPosts:     subPosts,
		Tags:         getResp.Tags,
	}

	setInfo(&resp.info, *getResp)

	SendResponse(c, nil, resp)
}
