package feed

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

type User struct {
	Name      string `json:"name"`
	Id        uint32 `json:"id"`
	AvatarUrl string `json:"avatar_url"`
} // @name User

type Source struct {
	// Kind uint32 `json:"kind"` // 类型，1 -> 团队，2 -> 项目，3 -> 文档，4 -> 文件，6 -> 进度（5 不使用）
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
} // @name Source

type FeedItem struct {
	Id          uint32  `json:"id"`
	Action      string  `json:"action"`
	ShowDivider bool    `json:"show_divider"` // 分割线
	CreateTime  string  `json:"create_time"`
	User        *User   `json:"user"`
	Source      *Source `json:"source"`
} // @name FeedItem

type FeedListResponse struct {
	List []*FeedItem `json:"list"`
} // @name FeedListResponse
