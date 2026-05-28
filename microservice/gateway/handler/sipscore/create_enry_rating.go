package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// CreateEntryRating ... 创建榜单条目评分
// @Summary 创建榜单条目评分 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateSipScoreEntryRatingRequest true "create_sip_score_entry_rating_request"
// @Success 200 {object} EmptyResponse
// @Router /sip-score/entry/rating [post]
func (a *Api) CreateEntryRating(c *gin.Context) {
	log.Info("Post CreateEntryRating function called.")

	var req CreateSipScoreEntryRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	createReq := &pb.CreateSipScoreEntryCommentRatingRequest{
		SipScoreId:      req.SipScoreID,
		SipScoreEntryId: req.EntryID,
		UserId:          userID,
		Comment:         req.Comment,
		Rating:          req.Rating,
		ImgUrl:          req.ImgUrl,
	}

	createResp, err := client.PostClient.CreateSipScoreEntryCommentRating(c.Request.Context(), createReq)
	if err != nil {
		SendError(c, err, createResp, "", GetLine())
		return
	}

	SendResponse(c, nil, &EmptyResponse{})
}
