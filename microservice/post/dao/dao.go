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
}

// Interface dao
type Interface interface {
	Create(*PostModel) error
	List(uint8) ([]*pb.Post, error)
	UpdateInfo(*PostModel) error
	Get(uint32) (*PostModel, error)
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
}

func GetDao() *Dao {
	return dao
}
