package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSipScoreEntry ... 获取榜单条目详情
// @Summary 获取榜单条目详情 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param sip_score_id path int true "sip_score_id"
// @Param entry_id path int true "entry_id"
// @Success 200 {object} GetSipScoreEntryResponse
// @Router /sip-score/entry/{sip_score_id}/{entry_id} [get]
func (a *Api) GetSipScoreEntry(c *gin.Context) {
	log.Info("SipScore GetSipScoreEntry function called.")

	sipScoreID, err := strconv.Atoi(c.Param("sip_score_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	entryID, err := strconv.Atoi(c.Param("entry_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	ok, err := model.Enforce(userID, constvar.SipScore, sipScoreID, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, &EmptyResponse{}, "权限不足", GetLine())
		return
	}

	detail, err := client.PostClient.GetSipScoreEntryDetail(c.Request.Context(), &pb.GetSipScoreEntryRequest{
		SipScoreId: uint32(sipScoreID),
		EntryId:    uint32(entryID),
		UserId:     userID,
	})
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, &GetSipScoreEntryResponse{
		Entry:    sipScoreEntryFromPB(detail.GetEntry()),
		MyRating: sipScoreEntryRatingFromPB(detail.GetMyRating()),
	})
}
