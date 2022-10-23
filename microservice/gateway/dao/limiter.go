package dao

import (
	"golang.org/x/time/rate"
	"time"
)

type LimiterManager struct {
	m map[uint32]*rate.Limiter
}

func initLimiterManager() *LimiterManager {
	m := make(map[uint32]*rate.Limiter)
	return &LimiterManager{
		m: m,
	}
}

func (l *LimiterManager) allowN(userId uint32, n int) bool {
	limiter, ok := l.m[userId]
	if !ok {
		limiter = rate.NewLimiter(1, 100) // 容量为100， 每秒产生0.1个 token
		l.m[userId] = limiter
	}

	return limiter.AllowN(time.Now(), n)
}
