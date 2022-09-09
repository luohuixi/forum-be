package post

import (
	"forum-gateway/dao"
	"forum-gateway/handler/comment"
	pb "forum-post/proto"
	"reflect"
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
	TypeName    string   `json:"type_name" binding:"required"` // normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal)
	Content     string   `json:"content" binding:"required"`
	Title       string   `json:"title,omitempty" binding:"required"`
	Category    string   `json:"category,omitempty" binding:"required"`
	ContentType string   `json:"content_type" binding:"required"` // md or rtf
	Tags        []string `json:"tags"`
}

type GetPostResponse struct {
	info
	Title        string     `json:"title"`
	Category     string     `json:"category"`
	IsCollection bool       `json:"is_collection"`
	SubPosts     []*SubPost `json:"sub_posts"`
	Tags         []string   `json:"tags"`
}

type SubPost struct {
	info
	Comments []*Comment `json:"comments"`
}

type Comment struct {
	info
	Replies []*info `json:"replies"`
}

type info struct {
	Id            uint32 `json:"id"`
	Content       string `json:"content"`
	CommentNum    uint32 `json:"comment_num"`
	Time          string `json:"time"`
	CreatorId     uint32 `json:"creator_id"`
	CreatorName   string `json:"creator_name"`
	CreatorAvatar string `json:"creator_avatar"`
	LikeNum       uint32 `json:"like_num"`
	IsLiked       bool   `json:"is_liked"`
}

func setInfo[T pb.CommentInfo | pb.Post](info *info, comment T) {
	typeT := reflect.TypeOf(comment)
	value := reflect.ValueOf(&comment).Elem()

	for i := 0; i < typeT.NumField(); i++ {
		v := value.Field(i)
		field := typeT.Field(i)

		switch field.Name {
		case "Id":
			info.Id = uint32(v.Uint())
		case "Content":
			info.Content = v.String()
		case "CommentNum":
			info.CommentNum = uint32(v.Uint())
		case "IsLiked":
			info.IsLiked = v.Bool()
		case "CreatorName":
			info.CreatorName = v.String()
		case "CreatorId":
			info.CreatorId = uint32(v.Uint())
		case "LikeNum":
			info.LikeNum = uint32(v.Uint())
		case "CreatorAvatar":
			info.CreatorAvatar = v.String()
		case "Time", "CreateTime":
			info.Time = v.String()
		}
	}
}
