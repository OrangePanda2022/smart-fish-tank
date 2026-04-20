package service

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	keyLimiterGCTick      = 10 * time.Minute
	keyLimiterExpireAfter = 30 * time.Minute
)

// KeyedRateLimiter 提供业务级别限流，例如登录/重置密码防爆破
type KeyedRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	lastSeen map[string]time.Time
	defaultR rate.Limit
	defaultB int
}

func NewKeyedRateLimiter(defaultR rate.Limit, defaultB int) *KeyedRateLimiter {
	// 每个业务 key 独立令牌桶，避免互相抢占配额
	l := &KeyedRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		lastSeen: make(map[string]time.Time),
		defaultR: defaultR,
		defaultB: defaultB,
	}
	go l.gc()
	return l
}

// Allow 使用默认速率规则判断是否放行
func (l *KeyedRateLimiter) Allow(key string) bool {
	return l.allowWith(key, l.defaultR, l.defaultB)
}

// allowWith 获取或创建指定 key 的限流器并消费一个令牌
func (l *KeyedRateLimiter) allowWith(key string, r rate.Limit, b int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter, ok := l.limiters[key]
	if !ok {
		limiter = rate.NewLimiter(r, b)
		l.limiters[key] = limiter
	}
	l.lastSeen[key] = time.Now()
	return limiter.Allow()
}

// gc 定期回收长时间未使用的 key，防止 map 持续增长
func (l *KeyedRateLimiter) gc() {
	ticker := time.NewTicker(keyLimiterGCTick)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		for key, seen := range l.lastSeen {
			if time.Since(seen) > keyLimiterExpireAfter {
				delete(l.lastSeen, key)
				delete(l.limiters, key)
			}
		}
		l.mu.Unlock()
	}
}
