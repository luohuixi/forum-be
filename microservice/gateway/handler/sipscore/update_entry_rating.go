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
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// UpdateEntryRating ... 修改评分
// @Summary 修改评分 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body UpdateSipScoreEntryRatingRequest true "update_sip_score_entry_rating_request"
// @Success 200 {object} handler.Response
// @Router /sip-score/entry/rating [put]
func (a *Api) UpdateEntryRating(c *gin.Context) {
	log.Info("SipScore UpdateEntryRating function called.")

	var req UpdateSipScoreEntryRatingRequest
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

	paths := buildRatingUpdatePaths(&req)
	if len(paths) == 0 {
		SendError(c, errno.ErrBadRequest, &EmptyResponse{}, "no fields to update", GetLine())
		return
	}

	updateReq := pb.UpdateSipScoreEntryCommentRatingInfoRequest{
		SipScoreId:      req.SipScoreID,
		SipScoreEntryId: req.EntryID,
		RatingId:        req.RatingID,
		LastModifiedBy:  userID,
		Rating:          req.Rating,
		ImgUrl:          req.ImgUrl,
		UpdateMask:      &fieldmaskpb.FieldMask{Paths: paths},
	}

	if req.Content != nil {
		updateReq.Content = *req.Content
	}

	_, err = client.PostClient.UpdateSipScoreEntryCommentRatingInfo(c.Request.Context(), &updateReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, &EmptyResponse{})
}

func buildRatingUpdatePaths(req *UpdateSipScoreEntryRatingRequest) []string {
	paths := make([]string, 0, 3)
	if req.Rating > 0 {
		paths = append(paths, "rating")
	}
	if req.Content != nil {
		paths = append(paths, "content")
	}
	if req.ImgUrl != "" {
		paths = append(paths, "img_url")
	}
	return paths
}
