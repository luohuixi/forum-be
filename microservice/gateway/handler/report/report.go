package report

import (
	"forum-gateway/dao"
)

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

type CreateRequest struct {
	TypeName string `json:"type_name" binding:"required"`
	Cause    string `json:"cause" binding:"required"`
	PostId   uint32 `json:"post_id" binding:"required"`
}

type HandleRequest struct {
	Id     uint32 `json:"id" binding:"required"`
	Result string `json:"result" binding:"required"` // invalid or valid
}

type ListResponse struct {
	Reports []*Report `json:"reports"`
}

type Report struct {
	Id                 uint32 `json:"id"`
	Cause              string `json:"cause"`
	PostId             uint32 `json:"post_id"`
	UserId             uint32 `json:"user_id"`
	TypeName           string `json:"type_name"`
	CreateTime         string `json:"create_time"`
	UserAvatar         string `json:"user_avatar"`
	UserName           string `json:"user_name"`
	BeReportedUserId   uint32 `json:"be_reported_user_id"`
	BeReportedUserName string `json:"be_reported_user_name"`
	PostTitle          string `json:"post_title"`
}
