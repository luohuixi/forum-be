package dao

import (
	"forum/model"
	"github.com/go-redis/redis"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	Redis *redis.Client
}

// Interface dao
type Interface interface {
	Post(queue string, data []byte) error
	CreateQueue(queue string) error
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	dao = &Dao{
		Redis: model.RedisDB.Self,
	}
}
