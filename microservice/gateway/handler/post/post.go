package post

import (
	"forum-gateway/dao"
	"forum-gateway/handler/comment"
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
	Id            uint32             `json:"id,omitempty"`
	Title         string             `json:"title,omitempty"`
	Time          string             `json:"time,omitempty"`
	Content       string             `json:"content,omitempty"`
	CategoryId    uint32             `json:"category_id,omitempty"`
	CreatorId     uint32             `json:"creator_id,omitempty"`
	CreatorName   string             `json:"creator_name,omitempty"`
	CreatorAvatar string             `json:"creator_avatar,omitempty"`
	CommentNum    uint32             `json:"comment_num,omitempty"`
	LikeNum       uint32             `json:"like_num,omitempty"`
	IsLiked       bool               `json:"is_liked,omitempty"`
	IsCollection  bool               `json:"is_collection,omitempty"`
	Comments      []*comment.Comment `json:"comments,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

type CreateRequest struct {
	TypeName   string   `json:"type_name"`
	Content    string   `json:"content"`
	Title      string   `json:"title,omitempty"`
	CategoryId uint32   `json:"category_id,omitempty"`
	MainPostId uint32   `json:"main_post_id"`
	Tags       []string `json:"tags"`
}
