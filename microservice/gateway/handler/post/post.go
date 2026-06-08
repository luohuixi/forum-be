package post

import (
	"forum-gateway/dao"
	"forum-gateway/handler/comment"
	pb "forum-post/proto"
)

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

// =====================
// Common
// =====================
type info struct {
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	CreateTime    string `json:"create_time"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	LikeNum       uint32 `json:"like_num"`
	IsLiked       bool   `json:"is_liked"`
}

// =====================
// Post Domain
// =====================

// ---- model ----

type Post struct {
	Id            uint32             `json:"id"`
	Title         string             `json:"title"`
	Time          string             `json:"time"`
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
	ContentType   string             `json:"content_type"` // md or rtf
	Summary       string             `json:"summary"`
	CollectionNum uint32             `json:"collection_num"`
}

type PostPartInfo struct {
	Id            uint32   `json:"id"`
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Category      string   `json:"category"`
	Time          string   `json:"time"`
	CreatorId     uint32   `json:"creator_id"`
	CreatorName   string   `json:"creator_name"`
	CreatorAvatar string   `json:"creator_avatar"`
	CommentNum    uint32   `json:"comment_num"`
	CollectionNum uint32   `json:"collection_num"`
	LikeNum       uint32   `json:"like_num"`
	Tags          []string `json:"tags"`
}

// ---- request ----

type CreateRequest struct {
	Domain          string   `json:"domain" binding:"required"` // normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal)
	Content         string   `json:"content" binding:"required"`
	CompiledContent string   `json:"compiled_content"`
	Title           string   `json:"title,omitempty" binding:"required"`
	Category        string   `json:"category,omitempty" binding:"required"`
	ContentType     string   `json:"content_type" binding:"required"` // md or rtf
	Summary         string   `json:"summary" binding:"required"`
	Tags            []string `json:"tags" binding:"required"`
}

type UpdateInfoRequest struct {
	Id       uint32   `json:"id" binding:"required"`
	Domain   string   `json:"domain" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags" binding:"required"`
}

type TrimHtmlRequest struct {
	Data string `json:"data" binding:"required"`
}

// ---- response ----

type GetPostResponse struct {
	info
	CommentNum      uint32     `json:"comment_num"`
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	Time            string     `json:"time"`
	IsCollection    bool       `json:"is_collection"`
	SubPosts        []*SubPost `json:"sub_posts"`
	Tags            []string   `json:"tags"`
	ContentType     string     `json:"content_type"` // md or rtf
	CompiledContent string     `json:"compiled_content"`
	Summary         string     `json:"summary"`
	CollectionNum   uint32     `json:"collection_num"`
}

type PostPartInfoResponse struct {
	Posts []*PostPartInfo `json:"posts"`
}

type ListMainPostResponse struct {
	Posts []*Post `json:"posts"`
}

// =====================
// Comment Domain
// =====================

// ---- model ----

type SubPost struct {
	info
	ImgUrl     string     `json:"img_url"`
	CommentNum uint32     `json:"comment_num"`
	Comments   []*Comment `json:"comments"`
}

type Comment struct {
	info
	BeRepliedId       uint32 `json:"be_replied_id"`
	BeRepliedUserId   uint32 `json:"be_replied_user_id"`
	BeRepliedContent  string `json:"be_replied_content"`
	BeRepliedUserName string `json:"be_replied_user_name"`
}

// =====================
// Other Common Response
// =====================

type QiNiuToken struct {
	Token string `json:"token"`
}

type UnReadNum struct {
	Num      uint32 `json:"num"`
	Category string `json:"category"`
}

type IdResponse struct {
	Id uint32 `json:"id"`
}

type UnReadNumResponse struct {
	UnReadNum []*UnReadNum `json:"un_read_num"`
}

// =====================
// internal util
// =====================

func setCommentInfo(info *info, comment *pb.CommentInfo) {
	info.Id = comment.Id
	info.Content = comment.Content
	info.CreateTime = comment.CreateTime.AsTime().Format("2006-01-02 15:04:05")
	info.CreatorId = comment.CreatorId
	info.CreatorName = comment.CreatorName
	info.CreatorAvatar = comment.CreatorAvatar
	info.LikeNum = comment.LikeNum
	info.IsLiked = comment.IsLiked
}

func setPostInfo(info *info, post *pb.Post) {
	info.Id = post.Id
	info.Content = post.Content
	info.CreateTime = post.Time
	info.CreatorId = post.CreatorId
	info.CreatorName = post.CreatorName
	info.CreatorAvatar = post.CreatorAvatar
	info.LikeNum = post.LikeNum
	info.IsLiked = post.IsLiked
}
