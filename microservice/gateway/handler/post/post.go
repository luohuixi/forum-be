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
	Id            uint32             `json:"id"`
	Title         string             `json:"title"`
	Time          string             `json:"time"`
	Content       string             `json:"content"`
	Category      string             `json:"category"`
	CreatorId     uint32             `json:"creator_id"`
	CreatorName   string             `json:"creator_name"`
	CreatorAvatar string             `json:"creator_avatar"`
	CommentNum    uint32             `json:"comment_num"`
	LikeNum       uint32             `json:"like_num"`
	IsLiked       bool               `json:"is_liked"`
	IsCollection  bool               `json:"is_collection"`
	Comments      []*comment.Comment `json:"comments"`
	Tags          []string           `json:"tags"`
}

type CreateRequest struct {
	TypeName   string   `json:"type_name" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Title      string   `json:"title,omitempty" binding:"required"`
	Category   string   `json:"category,omitempty" binding:"required"`
	MainPostId uint32   `json:"main_post_id"`
	Tags       []string `json:"tags"`
}
