package dao

import (
	"forum/model"
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

	PublishMsg([]byte) error
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

	// init redis pub-sub client
	model.PubSubClient.Init(RdbChan)

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}

const RdbChan = "sub"

func (d Dao) PublishMsg(msg []byte) error {
	return d.Redis.Publish(RdbChan, msg).Err()
}
