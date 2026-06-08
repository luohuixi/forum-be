package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	userpb "forum-user/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListEntryRatings ... 获取评分列表
// @Summary list 获取评分列表 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param sip_score_id path int true "sip_score_id"
// @Param entry_id path int true "entry_id"
// @Param sort_type query int false "sort_type"
// @Param page_size query int false "page_size"
// @Param page_token query string false "page_token"
// @Success 200 {object} ListSipScoreEntryRatingsResponse
// @Router /sip-score/entry-rating/list/{sip_score_id}/{entry_id} [get]
func (a *Api) ListEntryRatings(c *gin.Context) {
	log.Info("SipScore ListEntryRatings function called.",
		zap.String("X-Request-Id", c.GetString("X-Request-Id")))

	sipScoreID, err := strconv.Atoi(c.Param("sip_score_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if sipScoreID <= 0 {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, "sip_score_id must be positive", GetLine())
		return
	}

	entryID, err := strconv.Atoi(c.Param("entry_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if entryID <= 0 {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, "entry_id must be positive", GetLine())
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	sortType, err := strconv.Atoi(c.DefaultQuery("sort_type", "1"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
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

	listReq := &pb.ListSipScoreEntryCommentRatingInfoRequest{
		SipScoreId: uint32(sipScoreID),
		EntryId:    uint32(entryID),
		PageToken:  pageToken,
		PageSize:   uint32(pageSize),
		SortType:   uint32(sortType),
	}

	listResp, err := client.PostClient.ListSipScoreEntryCommentRating(c.Request.Context(), listReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	rpcRatings := listResp.Ratings

	// Collect all user IDs to batch fetch user info
	userIDSet := make(map[uint32]struct{})
	commentUserIDSet := make(map[uint32]struct{})
	for _, r := range rpcRatings {
		if r == nil {
			continue
		}
		userIDSet[r.CreatorId] = struct{}{}
		userIDSet[r.LastModifiedBy] = struct{}{}
		for _, cmt := range r.Comments {
			if cmt != nil {
				commentUserIDSet[cmt.CreatorId] = struct{}{}
			}
		}
	}
	for uid := range commentUserIDSet {
		userIDSet[uid] = struct{}{}
	}

	userIDs := make([]uint32, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	var userMap map[uint32]*userpb.UserInfo
	if len(userIDs) > 0 {
		userResp, err := client.UserClient.GetInfo(c.Request.Context(), &userpb.GetInfoRequest{Ids: userIDs})
		if err != nil {
			log.Error("Failed to get user info for ratings",
				zap.Uint32s("userIds", userIDs),
				zap.Error(err),
			)
		}
		if userResp != nil {
			userMap = make(map[uint32]*userpb.UserInfo, len(userResp.GetList()))
			for _, u := range userResp.GetList() {
				if u != nil {
					userMap[u.Id] = u
				}
			}
		}
	}

	buildUserInfo := func(id uint32) *userInfo {
		if u, ok := userMap[id]; ok {
			return &userInfo{ID: id, Name: u.Name, Avatar: u.AvatarUrl}
		}
		return &userInfo{ID: id}
	}

	isRatingLiked := func(ratingID uint32) bool {
		if ratingID == 0 {
			return false
		}
		key := "like:" + constvar.SipScoreEntryCommentRating + "_list:" + strconv.Itoa(int(ratingID))
		ok, err := model.SIsmembersFromRedis(key, userID)
		if err != nil {
			log.Error("Failed to get sip score rating like state",
				zap.Uint32("ratingId", ratingID),
				zap.Uint32("userId", userID),
				zap.Error(err),
			)
			return false
		}
		return ok
	}

	httpRatings := make([]*SipScoreEntryCommentRatingInfo, len(rpcRatings))
	for i, r := range rpcRatings {
		if r == nil {
			r = &pb.SipScoreEntryCommentRating{}
		}

		comments := make([]*CommentInfo, len(r.Comments))
		for j, cmt := range r.Comments {
			if cmt == nil {
				cmt = &pb.CommentInfo{}
			}
			comments[j] = &CommentInfo{
				ID:              cmt.Id,
				TypeName:        cmt.TypeName,
				Content:         cmt.Content,
				FatherID:        cmt.FatherId,
				CreateTime:      cmt.CreateTime.AsTime().Format(time.DateTime),
				CreatorID:       cmt.CreatorId,
				CreatorName:     cmt.CreatorName,
				CreatorAvatar:   cmt.CreatorAvatar,
				LikeNum:         cmt.LikeNum,
				IsLiked:         cmt.IsLiked,
				BeRepliedUserID: cmt.BeRepliedUserId,
				BeRepliedName:   cmt.BeRepliedUserName,
				FatherContent:   cmt.FatherContent,
				ImgUrl:          cmt.ImgUrl,
				TargetID:        cmt.TargetId,
				TargetType:      cmt.TargetType,
			}
		}

		httpRatings[i] = &SipScoreEntryCommentRatingInfo{
			ID:              r.Id,
			SipScoreID:      r.SipScoreId,
			SipScoreEntryID: r.SipScoreEntryId,
			Creator:         buildUserInfo(r.CreatorId),
			LastModifiedBy:  buildUserInfo(r.LastModifiedBy),
			Rating:          r.Rating,
			Content:         r.Content,
			CommentID:       r.CommentId,
			LikeNum:         r.LikeNum,
			IsLiked:         isRatingLiked(r.Id),
			ImgUrl:          r.ImgUrl,
			CreatedAt:       r.CreatedAt.AsTime().Format(time.DateTime),
			UpdatedAt:       r.UpdatedAt.AsTime().Format(time.DateTime),
			CommentNum:      r.CommentNum,
			Comments:        comments,
		}
	}

	resp := &ListSipScoreEntryRatingsResponse{
		Ratings:   httpRatings,
		PageToken: listResp.PageToken,
		HasMore:   listResp.HasMore,
	}

	SendResponse(c, nil, resp)
}
