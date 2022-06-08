package dao

import (
	"forum/model"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
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
	GetList(string, time.Duration) ([]string, error)
	Rewrite(string, []string) error
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
