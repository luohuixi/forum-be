package dao

import (
	"forum-chat/service"
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
	Create(*service.ChatData) error
	GetList(uint32) ([]string, error)
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
