package user

import (
	"context"
	. "forum-gateway/handler"
	"strconv"

	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetProfile ... 获取 userProfile
// @Summary 获取用户 profile api
// @Description 通过 userId 获取完整 user 信息（权限: Normal-普通学生用户; NormalAdmin-学生管理员; Muxi-团队成员; MuxiAdmin-团队管理员; SuperAdmin-超级管理员）
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param id path int true "user_id"
// @Success 200 {object} UserProfile
// @Router /user/profile/{id} [get]
func GetProfile(c *gin.Context) {
	log.Info("User GetProfile function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	user, err := GetUserProfile(uint32(id))

	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, user)
}

func GetUserProfile(id uint32) (*UserProfile, error) {

	getProfileReq := &pb.GetRequest{Id: id}

	getProfileResp, err := service.UserClient.GetProfile(context.TODO(), getProfileReq)
	if err != nil {
		return nil, err
	}

	// 构造返回 response
	resp := &UserProfile{
		Id:        getProfileResp.Id,
		Name:      getProfileResp.Name,
		Avatar:    getProfileResp.Avatar,
		Email:     getProfileResp.Email,
		Role:      getProfileResp.Role,
		Signature: getProfileResp.Signature,
	}

	return resp, nil
}
