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
	Id       uint32   `json:"id" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags" binding:"required"`
}

type Post struct {
	Id            uint32             `json:"id,omitempty"`
	Title         string             `json:"title,omitempty"`
	Time          string             `json:"time,omitempty"`
	Content       string             `json:"content,omitempty"`
	Category      string             `json:"category,omitempty"`
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
	TypeName   string   `json:"type_name" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Title      string   `json:"title,omitempty" binding:"required"`
	Category   string   `json:"category,omitempty" binding:"required"`
	MainPostId uint32   `json:"main_post_id"`
	Tags       []string `json:"tags"`
}
