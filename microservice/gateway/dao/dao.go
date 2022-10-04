package dao

import (
	// pb "forum-post/proto"
	"forum/model"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	// Redis *redis.Client
}

// Interface dao
type Interface interface{}

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
		// Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}
