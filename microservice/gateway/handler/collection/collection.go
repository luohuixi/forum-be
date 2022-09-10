package collection

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
	PostId uint32 `json:"post_id,omitempty" binding:"required"`
}

type Collection struct {
	Id            uint32 `json:"id"`
	PostId        uint32 `json:"post_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Time          string `json:"time"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	CommentNum    uint32 `json:"comment_num"`
}
