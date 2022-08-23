package dao

import (
	"forum/model"
	m "forum/model"
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
}

// Interface dao
type Interface interface {
	Create(*FeedModel) (uint32, error)
	Delete(uint32) error
	List(*FeedModel, uint32, uint32, uint32, bool) ([]*FeedModel, error)
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
	// model.CB.Init()

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}

const RdbChan = "sub"

func PublishMsg(msg []byte) error {
	return m.RedisDB.Self.Publish(RdbChan, msg).Err()
}
