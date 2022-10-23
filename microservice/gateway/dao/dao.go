package dao

import (
	"forum/model"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	LimiterManager *LimiterManager
	// Redis *redis.Client
}

// Interface dao
type Interface interface {
	AllowN(userId uint32, n int) bool
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

	limiterManager := initLimiterManager()

	dao = &Dao{
		LimiterManager: limiterManager,
		// Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}

func (d Dao) AllowN(userId uint32, n int) bool {
	return d.LimiterManager.allowN(userId, n)
}
