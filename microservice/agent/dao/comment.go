package dao

import (
	"forum/pkg/constvar"
	"forum/util"
)

type CommentModel struct {
	Id         uint32
	TypeName   string // constvar.SubPost
	Content    string
	FatherId   uint32
	CreateTime string
	Re         bool
	CreatorId  uint32
	PostId     uint32
	LikeNum    uint32
	ImgUrl     string
	IsReport   bool
}

type CommentAgentReturn struct {
	Content string `json:"content"`
	PostId  uint32 `json:"post_id"`
}

func (c *CommentModel) Create() error {
	return dao.DB.Table("comments").Create(c).Error
}

func ChangeToCommentModel(c *CommentAgentReturn) *CommentModel {
	return &CommentModel{
		TypeName:   constvar.SubPost,
		Content:    c.Content,
		FatherId:   c.PostId,
		CreateTime: util.GetCurrentTime(),
		Re:         false,
		CreatorId:  163, // 暂定是root
		PostId:     c.PostId,
		LikeNum:    0,
		ImgUrl:     "",
		IsReport:   false,
	}
}
