package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetInfo ... 获取 userInfo
func GetInfo(c *gin.Context) {
	log.Info("User getInfo function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	// 从前端获取 Ids
	var req GetInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if len(req.Ids) == 0 {
		SendResponse(c, nil, &GetInfoResponse{})
		return
	}

	// 构造请求给 getInfo
	getInfoReq := &pb.GetInfoRequest{}
	getInfoReq.Ids = make([]uint32, len(req.Ids))
	for i, id := range req.Ids {
		getInfoReq.Ids[i] = id
	}

	getInfoResp, err := service.UserClient.GetInfo(context.TODO(), getInfoReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 构造返回 response
	var resp GetInfoResponse
	for _, item := range getInfoResp.List {
		resp.List = append(resp.List, userInfo{
			Id:        item.Id,
			Name:      item.Name,
			AvatarURL: item.AvatarUrl,
			Email:     item.Email,
		})
	}

	SendResponse(c, nil, resp)
}
