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

// -------------
// model
// -------------

type Comment struct {
	Id               uint32 `json:"id"`
	Content          string `json:"content"`
	TypeName         string `json:"type_name"` // first-level -> 一级评论; second-level -> 其它级
	FatherId         uint32 `json:"father_id"`
	CreateTime       string `json:"create_time"`
	CreatorId        uint32 `json:"creator_id"`
	CreatorName      string `json:"creator_name"`
	CreatorAvatar    string `json:"creator_avatar"`
	LikeNum          uint32 `json:"like_num"`
	IsLiked          bool   `json:"is_liked"`
	BeRepliedUserId  uint32 `json:"be_replied_user_id"`
	BeRepliedContent string `json:"be_replied_content"`
	ImgUrl           string `json:"img_url"`
}

// CommentItem 匹配 proto CommentInfo 的 JSON 输出
type CommentItem struct {
	Id                uint32         `json:"id"`
	TypeName          string         `json:"type_name"`
	Content           string         `json:"content"`
	FatherId          uint32         `json:"father_id"`
	CreateTime        string         `json:"create_time"`
	CreatorId         uint32         `json:"creator_id"`
	CreatorName       string         `json:"creator_name"`
	CreatorAvatar     string         `json:"creator_avatar"`
	LikeNum           uint32         `json:"like_num"`
	IsLiked           bool           `json:"is_liked"`
	BeRepliedUserId   uint32         `json:"be_replied_user_id"`
	BeRepliedUserName string         `json:"be_replied_user_name"`
	FatherContent     string         `json:"father_content"`
	ImgUrl            string         `json:"img_url"`
	TargetId          uint32         `json:"target_id"`
	TargetType        string         `json:"target_type"`
	SubNum            uint32         `json:"sub_num"`
	SubComments       []*CommentItem `json:"sub_comments"`
}

// ----------
// request
// ----------

type CreateRequest struct {
	TargetId   uint32 `json:"target_id" binding:"required"`
	TargetType string `json:"target_type" binding:"required"`
	TypeName   string `json:"type_name" binding:"required"`
	FatherId   uint32 `json:"father_id"`
	Content    string `json:"content" binding:"required"`
	ImgUrl     string `json:"img_url"`
}

type ListRequest struct {
	TargetId   uint32 `json:"target_id"`
	TargetType string `json:"target_type"`
	PageToken  string `json:"page_token"`
	PageSize   uint32 `json:"page_size"`
	SortType   uint32 `json:"sort_type"`
	FatherId   uint32 `json:"father_id"` // 可选：按父评论 ID 列出子评论（二级评论分页）
}

// ----------
// response
// ----------

type ListResponse struct {
	Comments  []*CommentItem `json:"comments"`
	PageToken string         `json:"page_token"`
	HasMore   bool           `json:"has_more"`
}
