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
	Create(*ChatData) error
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

// ChatData 发送到redis里面的数据
type ChatData struct {
	Message  string `json:"message"`
	Date     string `json:"date"`
	Receiver string `json:"receiver"`
	Sender   uint32 `json:"sender"`
}

func GetDao() *Dao {
	return dao
}
