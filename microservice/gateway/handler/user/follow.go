package user

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Follow ... 关注/取关用户
// @Summary 关注/取关用户 api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body FollowRequest true "follow_request"
// @Success 200 {object} FollowResponse
// @Router /user/follow [post]
func Follow(c *gin.Context) {
	log.Info("User Follow function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	resp, err := client.UserClient.ToggleFollow(c.Request.Context(), &pb.FollowRequest{
		UserId:       userID,
		TargetUserId: req.TargetUserID,
	})
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, resp, FollowResponse{})
}
