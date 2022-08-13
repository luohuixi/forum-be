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
	TypeName string `json:"type_name"`
}

type ListResponse struct {
	likes *[]Item
}
