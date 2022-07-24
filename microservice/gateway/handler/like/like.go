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
	TargetId uint32 `json:"target_id"`
	TypeId   uint32 `json:"type_id"`
}

type ListResponse struct {
	likes *[]Item
}
