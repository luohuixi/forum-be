package dao

import (
	pb "forum-post/proto"
	"forum/model"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB    *gorm.DB
	Redis *redis.Client
	CB    *casbin.Enforcer
	Map   map[uint8]string // TypeId2Name
}

// Interface dao
type Interface interface {
	CreatePost(*PostModel) error
	ListPost(uint8) ([]*PostInfo, error)
	UpdatePostInfo(*PostModel) error
	GetPost(uint32) (*PostModel, error)
	GetPostInfo(uint32) (*PostInfo, error)

	CreateComment(*CommentModel) error
	GetComment(uint32) (*CommentModel, error)
	UpdateCommentInfo(*CommentModel) error
	ListCommentByPostId(uint32) ([]*pb.CommentInfo, error)
	GetCommentNumByPostId(uint32) uint32

	AddLike(uint32, Item) error
	RemoveLike(uint32, Item) error
	GetLikedNum(Item) (int64, error)
	IsUserHadLike(uint32, Item) (bool, error)
	ListUserLike(uint32) ([]*Item, error)

	Enforce(...interface{}) (bool, error)
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// init db
	model.DB.Init()

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
		CB:    model.CB.Self,
	}

	dao.Map[1] = "post"
	dao.Map[2] = "comment"
}

func GetDao() *Dao {
	return dao
}

func (d *Dao) Enforce(rvals ...interface{}) (bool, error) {
	return d.CB.Enforce(rvals)
}
