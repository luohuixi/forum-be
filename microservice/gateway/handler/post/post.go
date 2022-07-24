package post

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

type UpdateInfoRequest struct {
	Id       uint32 `json:"id"`
	Content  string `json:"content"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type Post struct {
	Content       string `json:"content,omitempty"`
	Title         string `json:"title,omitempty"`
	LastEditTime  string `json:"last_edit_time,omitempty"`
	Category      string `json:"category,omitempty"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name,omitempty"`
	CreatorAvatar string `json:"creator_avatar"`
	IsLiked       bool   `json:"is_liked"`
	IsFavorite    bool   `json:"is_favorite"`
}

type ListResponse struct {
	posts *[]Post
}

type CreateRequest struct {
	TypeId     uint32 `json:"type_id"`
	Content    string `json:"content"`
	Title      string `json:"title,omitempty"`
	Category   string `json:"category,omitempty"`
	MainPostId uint32 `json:"main_post_id"`
}

type ListMainPostRequest struct {
	TypeId   uint32 `json:"type_id"`
	Category string `json:"category,omitempty"`
}

type ListSubPostRequest struct {
	TypeId     uint32 `json:"type_id"`
	MainPostId uint32 `json:"main_post_id"`
}
