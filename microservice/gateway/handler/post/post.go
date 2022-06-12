package post

import "forum-gateway/dao"

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

type UpdateInfoRequest struct {
	Id       uint32 `json:"id"`
	Content  string
	Title    string
	Category string
}
