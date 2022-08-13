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
	Id         uint32 `json:"id"`
	Content    string `json:"content"`
	Title      string `json:"title"`
	CategoryId uint32 `json:"category_id"`
}

type Post struct {
	Content       string `json:"content,omitempty"`
	Title         string `json:"title,omitempty"`
	LastEditTime  string `json:"last_edit_time,omitempty"`
	CategoryId    uint32 `json:"category_id,omitempty"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name,omitempty"`
	CreatorAvatar string `json:"creator_avatar"`
	IsLiked       bool   `json:"is_liked"`
	IsFavorite    bool   `json:"is_favorite"`
}

type ListResponse struct {
	Posts []Post `json:"posts"`
}

type CreateRequest struct {
	TypeName   string `json:"type_name"`
	Content    string `json:"content"`
	Title      string `json:"title,omitempty"`
	CategoryId uint32 `json:"category_id,omitempty"`
	MainPostId uint32 `json:"main_post_id"`
}
