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
	"time"

	"github.com/gin-gonic/gin"
)

// SearchEntries ... 搜索榜单条目
// @Summary search 搜索榜单条目 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param sip_score_id path int true "sip_score_id"
// @Param keyword query string true "keyword"
// @Param page_size query int false "page_size"
// @Param page_token query string false "page_token"
// @Success 200 {object} ListSipScoreEntriesResponse
// @Router /sip-score/entries/search/{sip_score_id} [get]
func (a *Api) SearchEntries(c *gin.Context) {
	log.Info("Post SearchEntries function called.")

	sipScoreID, err := strconv.Atoi(c.Param("sip_score_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, "keyword required", GetLine())
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pageToken := c.DefaultQuery("page_token", "")

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
	if ok = a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, &EmptyResponse{}, "Please try again later", GetLine())
		return
	}

	searchReq := &pb.SearchSipScoreEntryRequest{
		SipScoreId: uint32(sipScoreID),
		Keyword:    keyword,
		PageToken:  pageToken,
		PageSize:   uint32(pageSize),
	}

	searchResp, err := client.PostClient.SearchSipScoreEntry(c.Request.Context(), searchReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	rpcEntries := searchResp.Entries
	httpEntries := make([]*SipScoreEntry, len(rpcEntries))
	for i, rpcEntry := range rpcEntries {
		if rpcEntry == nil {
			rpcEntry = &pb.SipScoreEntry{}
		}
		httpEntries[i] = &SipScoreEntry{
			ID:               rpcEntry.Id,
			SipScoreID:       rpcEntry.SipScoreId,
			CreatedAt:        rpcEntry.CreatedAt.AsTime().Format(time.DateTime),
			UpdatedAt:        rpcEntry.UpdatedAt.AsTime().Format(time.DateTime),
			Creator:          &userInfo{ID: rpcEntry.CreatorId},
			LastModifiedBy:   &userInfo{ID: rpcEntry.LastModifiedBy},
			Name:             rpcEntry.Name,
			Description:      rpcEntry.Description,
			CoverImg:         rpcEntry.CoverImg,
			ParticipantCount: rpcEntry.ParticipantCount,
			CommentCount:     rpcEntry.CommentCount,
			ScoreTotal:       rpcEntry.ScoreTotal,
			ScoreAvg:         rpcEntry.ScoreAvg,
		}
	}

	resp := &ListSipScoreEntriesResponse{
		Entries:   httpEntries,
		PageToken: searchResp.PageToken,
		HasMore:   searchResp.HasMore,
	}

	SendResponse(c, nil, resp)
}
