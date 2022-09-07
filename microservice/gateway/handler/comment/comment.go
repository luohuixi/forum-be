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
	TypeName string `json:"type_name,omitempty" binding:"required"`
	Content  string `json:"content,omitempty" binding:"required"`
	FatherId uint32 `json:"father_id,omitempty" binding:"required"`
	PostId   uint32 `json:"post_id,omitempty" binding:"required"`
}

type Comment struct {
	Id            uint32 `json:"id,omitempty"`
	Content       string `json:"content,omitempty"`
	TypeName      string `json:"type_name,omitempty"`
	FatherId      uint32 `json:"father_id,omitempty"`
	CreateTime    string `json:"create_time,omitempty"`
	CreatorId     uint32 `json:"creator_id,omitempty"`
	CreatorName   string `json:"creator_name,omitempty"`
	CreatorAvatar string `json:"creator_avatar,omitempty"`
	LikeNum       uint32 `json:"like_num,omitempty"`
	IsLiked       bool   `json:"is_liked,omitempty"`
}
