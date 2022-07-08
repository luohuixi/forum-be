package dao

import (
	// pb "forum-post/proto"
	"forum/model"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	Redis *redis.Client
	CB    *casbin.Enforcer
}

// Interface dao
type Interface interface {
	Enforce(...interface{}) (bool, error)
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// 黑名单过期数据定时清理
	// go service.TidyBlacklist()
	// 同步黑名单数据
	// service.SynchronizeBlacklistToRedis()

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	dao = &Dao{
		Redis: model.RedisDB.Self,
		CB:    model.CB.Self,
	}
}

func GetDao() *Dao {
	return dao
}

func (d *Dao) Enforce(rvals ...interface{}) (bool, error) {
	return d.CB.Enforce(rvals)
}
