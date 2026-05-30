package comment

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... 批量获取评论列表
// @Summary 批量获取评论列表 api
// @Tags comment
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body ListRequest true "list_comment_request"
// @Success 200 {object} ListResponse
// @Router /comment/list [post]
func (a *Api) List(c *gin.Context) {
	log.Info("Comment List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req ListRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	// 按 father_id 列出子评论时不需要校验 target_id + target_type
	if req.FatherId == 0 {
		if req.TargetId == 0 || req.TargetType == "" {
			SendError(c, errno.ErrBadRequest, nil, "target_id and target_type required", GetLine())
			return
		}
	}

	if req.SortType != constvar.SortByNewest && req.SortType != constvar.SortByHottest {
		SendError(c, errno.ErrBadRequest, nil, "sort_type not legal", GetLine())
		return
	}

	// TODO: 添加权限校验，分别对 father_id 和 target_id+target_type 两种场景进行鉴权
	_ = userId

	listReq := &pb.ListCommentRequest{
		TargetId:   req.TargetId,
		TargetType: req.TargetType,
		PageToken:  req.PageToken,
		PageSize:   req.PageSize,
		UserId:     userId,
		SortType:   req.SortType,
		FatherId:   req.FatherId,
	}

	listResp, err := client.PostClient.ListComments(c.Request.Context(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, ListResponse{})
}
