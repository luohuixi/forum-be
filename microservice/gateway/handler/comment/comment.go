package comment

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
	TypeId   uint32
	Content  string
	FatherId uint32
	PostId   uint32
}
