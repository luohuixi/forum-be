package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// DeleteEntryRating ... 删除评分
// @Summary 删除评分 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body DeleteSipScoreEntryRatingRequest true "delete_sip_score_entry_rating_request"
// @Success 200 {object} Response
// @Router /sip-score/entry/rating [delete]
func (a *Api) DeleteEntryRating(c *gin.Context) {
	log.Info("SipScore DeleteEntryRating function called.")

	var req DeleteSipScoreEntryRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	// 权限检查：创建者拥有该 rating 的 Write 权限，管理员通过 g2 资源域匹配也能通过
	ok, err := model.Enforce(userID, constvar.SipScoreEntryCommentRating, req.RatingID, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, &EmptyResponse{}, "权限不足", GetLine())
		return
	}
	if ok = a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, &EmptyResponse{}, "Please try again later", GetLine())
		return
	}

	deleteReq := &pb.DeleteSipScoreEntryCommentRatingRequest{
		SipScoreId:      req.SipScoreID,
		SipScoreEntryId: req.EntryID,
		RatingId:        req.RatingID,
		LastModifiedBy:  userID,
	}

	_, err = client.PostClient.DeleteSipScoreEntryCommentRating(c.Request.Context(), deleteReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, &EmptyResponse{})
}
