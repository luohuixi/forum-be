package agent

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

type AddKnowledgeRequest struct {
	PostId uint32 `json:"post_id" binding:"required"`
}

type GiveAnswerRequest struct {
	PostId       uint32 `json:"post_id" binding:"required"`
	ExtraContent string `json:"extra_content" binding:"required"`
}
