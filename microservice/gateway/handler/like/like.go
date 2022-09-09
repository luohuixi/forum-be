package like

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

type Item struct {
	TargetId uint32 `json:"target_id" binding:"required"`
	TypeName string `json:"type_name" binding:"required"` // post or comment
}

type ListResponse struct {
	likes *[]Item
}
