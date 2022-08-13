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
	Time     string `json:"time"`
	Receiver string `json:"-"`
	Sender   string `json:"sender"`
	TypeName string `json:"type_name"`
}

func GetDao() *Dao {
	return dao
}
