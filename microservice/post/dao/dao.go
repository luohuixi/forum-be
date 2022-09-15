package dao

import (
	"errors"
	pb "forum-post/proto"
	"forum/model"
	"forum/pkg/constvar"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// Interface dao
type Interface interface {
	DeleteItem(Item) error

	CreatePost(*PostModel) (uint32, error)
	ListPost(*PostModel, uint32, uint32, uint32, bool, string) ([]*PostInfo, error)
	GetPostInfo(uint32) (*PostInfo, error)
	GetPost(uint32) (*PostModel, error)
	IsUserCollectionPost(uint32, uint32) (bool, error)

	CreateComment(*CommentModel) (uint32, error)
	GetCommentInfo(uint32) (*CommentInfo, error)
	GetComment(uint32) (*CommentModel, error)
	ListCommentByPostId(uint32) ([]*pb.CommentInfo, error)
	GetCommentNumByPostId(uint32) (uint32, error)

	AddLike(uint32, Item) error
	RemoveLike(uint32, Item) error
	GetLikedNum(Item) (int64, error)
	IsUserHadLike(uint32, Item) (bool, error)
	ListUserLike(uint32) ([]*Item, error)

	CreatePost2Tag(Post2TagModel) error
	GetTagById(uint32) (*TagModel, error)
	GetTagByContent(string) (*TagModel, error)
	ListTagsByPostId(uint32) ([]string, error)

	AddTagToSortedSet(uint32) error
	ListPopularTags() ([]string, error)

	CreateCollection(*CollectionModel) (uint32, error)
	DeleteCollection(*CollectionModel) error
	ListCollectionByUserId(uint32) ([]*pb.Collection, error)
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
	}
}

func GetDao() *Dao {
	return dao
}

func (Dao) DeleteItem(i Item) error {
	if i.TypeName == constvar.Post {
		item := &PostModel{}
		if err := item.Get(i.Id); err != nil {
			return err
		}
		return item.Delete()
	} else if i.TypeName == constvar.Comment {
		item := &CommentModel{}
		if err := item.Get(i.Id); err != nil {
			return err
		}
		return item.Delete()
	} else {
		return errors.New("wrong TypeName")
	}
}
