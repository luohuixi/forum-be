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
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags" binding:"required"`
}

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

type CreateRequest struct {
	TypeName        string   `json:"type_name" binding:"required"` // normal -> 团队外; muxi -> 团队内 (type_name暂时均填normal)
	Content         string   `json:"content" binding:"required"`
	CompiledContent string   `json:"compiled_content"`
	Title           string   `json:"title,omitempty" binding:"required"`
	Category        string   `json:"category,omitempty" binding:"required"`
	ContentType     string   `json:"content_type" binding:"required"` // md or rtf
	Summary         string   `json:"summary" binding:"required"`
	Tags            []string `json:"tags" binding:"required"`
}

type GetPostResponse struct {
	info
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	IsCollection    bool       `json:"is_collection"`
	SubPosts        []*SubPost `json:"sub_posts"`
	Tags            []string   `json:"tags"`
	ContentType     string     `json:"content_type"` // md or rtf
	CompiledContent string     `json:"compiled_content"`
	Summary         string     `json:"summary"`
	CollectionNum   uint32     `json:"collection_num"`
}

type SubPost struct {
	info
	Comments []*Comment `json:"comments"`
}

type Comment struct {
	info
	BeRepliedId      uint32 `json:"be_replied_id"`
	BeRepliedContent string `json:"be_replied_content"`
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

type PostPartInfo struct {
	PostId        uint32   `json:"post_id"`
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

type PostPartInfoResponse struct {
	Posts []*PostPartInfo `json:"posts"`
}

type ListMainPostResponse struct {
	Posts []*Post `json:"posts"`
}

// 这里用了 generics and reflect, 更一般的写法应该是用interface
func setInfo[T pb.CommentInfo | pb.Post](info *info, comment *T) {
	typeT := reflect.TypeOf(*comment)
	value := reflect.ValueOf(comment).Elem()

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

type QiNiuToken struct {
	Token string `json:"token"`
}

func GetPostPartInfoResponse(data *pb.ListPostPartInfoResponse) *PostPartInfoResponse {
	posts := make([]*PostPartInfo, len(data.Posts))
	for i, p := range data.Posts {
		posts[i] = &PostPartInfo{
			PostId:        p.PostId,
			Title:         p.Title,
			Summary:       p.Summary,
			Category:      p.Category,
			Time:          p.Time,
			CreatorId:     p.CreatorId,
			CreatorName:   p.CreatorName,
			CreatorAvatar: p.CreatorAvatar,
			CommentNum:    p.CommentNum,
			LikeNum:       p.LikeNum,
			Tags:          p.Tags,
			CollectionNum: p.CollectionNum,
		}
	}

	return &PostPartInfoResponse{
		Posts: posts,
	}
}
