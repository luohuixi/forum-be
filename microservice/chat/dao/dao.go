package dao

import (
	"forum/model"
	"github.com/go-redis/redis"
	"time"
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
	GetList(string, time.Duration) ([]string, error)
	Rewrite(string, []string) error
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// init redis
	model.RedisDB.Init()

	dao = &Dao{
		Redis: model.RedisDB.Self,
	}

}

// ChatData 发送到redis里面的数据
type ChatData struct {
	Content  string `json:"content"`
	Date     string `json:"date"`
	Receiver string `json:"-"`
	Sender   string `json:"sender"`
	TypeId   uint32 `json:"type_id"`
}

func GetDao() *Dao {
	return dao
}
