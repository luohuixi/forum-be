package sipscore

import (
	. "forum-gateway/handler"
	"forum-gateway/handler/user"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *Api) listUserSipScores(c *gin.Context, collected bool) {
	targetUserID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	viewerID := c.MustGet("userId").(uint32)
	if collected && int(viewerID) != targetUserID {
		ok, err := model.Enforce(viewerID, constvar.CollectionAndLike, targetUserID, constvar.Read)
		if err != nil {
			SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
			return
		}
		if !ok {
			public, err := user.IsCollectionAndLikePublic(c.Request.Context(), uint32(targetUserID))
			if err != nil {
				SendError(c, err, &EmptyResponse{}, "", GetLine())
				return
			}
			ok = public
		}
		if !ok {
			SendResponse(c, errno.ErrPrivacyInfo, nil)
			return
		}
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	lastID, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	req := &pb.ListPostPartInfoRequest{
		UserId:       viewerID,
		TargetUserId: uint32(targetUserID),
		LastId:       uint32(lastID),
		Offset:       uint32(page * limit),
		Limit:        uint32(limit),
		Pagination:   limit != 0 || page != 0,
	}

	var resp *pb.ListSipScoreResponse
	if collected {
		resp, err = client.PostClient.ListUserCollectedSipScore(c.Request.Context(), req)
	} else {
		resp, err = client.PostClient.ListUserCreatedSipScore(c.Request.Context(), req)
	}
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, listSipScoresResponseFromPB(resp))
}

// ListUserCreatedSipScores ... 获取用户创建的茶评榜单
// @Router /sip-score/created/{user_id} [get]
func (a *Api) ListUserCreatedSipScores(c *gin.Context) {
	log.Info("SipScore ListUserCreatedSipScores function called.")
	a.listUserSipScores(c, false)
}

// ListUserCollectedSipScores ... 获取用户收藏的茶评榜单
// @Router /sip-score/collected/{user_id} [get]
func (a *Api) ListUserCollectedSipScores(c *gin.Context) {
	log.Info("SipScore ListUserCollectedSipScores function called.")
	a.listUserSipScores(c, true)
}
