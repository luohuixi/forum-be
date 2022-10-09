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
	TypeName string `json:"type_name" binding:"required"` // sub-post -> 从帖; first-level -> 一级评论; second-level -> 其它级
	Content  string `json:"content" binding:"required"`
	FatherId uint32 `json:"father_id" binding:"required"`
	PostId   uint32 `json:"post_id" binding:"required"`
}

type Comment struct {
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	TypeName      string `json:"type_name"` // first-level -> 一级评论; second-level -> 其它级
	FatherId      uint32 `json:"father_id"`
	CreateTime    string `json:"create_time"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	LikeNum       uint32 `json:"like_num"`
	IsLiked       bool   `json:"is_liked"`
}
