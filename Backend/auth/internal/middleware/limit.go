package middleware

import (
	errorHandler "auth/internal/error"
	"context"
	"hash/fnv"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

const (
	shardCount           = 32
	maxLimiterPerShard   = 10000
	localGCTick          = 5 * time.Minute
	localExpireWindow    = 15 * time.Minute
	localGCShardInterval = 10 * time.Millisecond
	redisLimiterTimeout  = 500 * time.Millisecond
	redisWindowMillis    = 1000
)

type ipLimiter struct {
	shards []*shard
	rate   rate.Limit
	burst  int
}

type shard struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	lastSeen map[string]time.Time
}

func newIPLimiter(r rate.Limit, b int, logger *slog.Logger) *ipLimiter {
	_ = logger

	l := &ipLimiter{
		shards: make([]*shard, shardCount),
		rate:   r,
		burst:  b,
	}
	for i := 0; i < shardCount; i++ {
		l.shards[i] = &shard{
			limiters: make(map[string]*rate.Limiter),
			lastSeen: make(map[string]time.Time),
		}
	}
	// 启动全局 GC
	go l.gc()
	return l
}

// getShard 根据 IP 哈希定位到对应分片
func (l *ipLimiter) getShard(ip string) *shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(ip))
	return l.shards[h.Sum32()%shardCount]
}

func (l *ipLimiter) Get(ip string) *rate.Limiter {
	s := l.getShard(ip)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果单分片容量过大拒绝创建
	if _, ok := s.limiters[ip]; !ok && len(s.limiters) > maxLimiterPerShard {
		// 返回一个全局兜底限流器，防止 OOM
		return rate.NewLimiter(l.rate, l.burst)
	}

	limiter, ok := s.limiters[ip]
	if !ok {
		limiter = rate.NewLimiter(l.rate, l.burst)
		s.limiters[ip] = limiter
	}
	s.lastSeen[ip] = time.Now()
	return limiter
}

func (l *ipLimiter) gc() {
	ticker := time.NewTicker(localGCTick)
	defer ticker.Stop()

	for range ticker.C {
		for _, s := range l.shards {
			s.mu.Lock()
			for ip, seen := range s.lastSeen {
				if time.Since(seen) > localExpireWindow {
					delete(s.lastSeen, ip)
					delete(s.limiters, ip)
				}
			}
			s.mu.Unlock()
			// 避免长时间占用 CPU
			time.Sleep(localGCShardInterval)
		}
	}
}

// HTTPRateLimit 提供接口级限流，默认按客户端 IP 控制请求速率
func HTTPRateLimitLocal(r rate.Limit, b int, logger *slog.Logger) gin.HandlerFunc {
	store := newIPLimiter(r, b, logger)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !store.Get(ip).Allow() {
			AbortWithError(c, errorHandler.ErrTooManyRequests)
			return
		}
		c.Next()
	}
}

type redisLimiter struct {
	client redis.UniversalClient
	burst  int
	cb     *gobreaker.CircuitBreaker
	logger *slog.Logger
}

func NewRedisLimiter(client redis.UniversalClient, burst int, logger *slog.Logger) *redisLimiter {
	// 配置熔断器
	st := gobreaker.Settings{
		Name:        "Redis-Limiter",
		MaxRequests: 5,                // 半开状态下允许通过的请求数
		Interval:    10 * time.Second, // 清空失败计数的周期
		Timeout:     5 * time.Second,  // 熔断器开启后多久尝试恢复 进入半开状态
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 如果连续失败超过 5 次，触发熔断
			return counts.ConsecutiveFailures > 5
		},
	}

	return &redisLimiter{
		client: client,
		burst:  burst,
		cb:     gobreaker.NewCircuitBreaker(st),
		logger: logger,
	}
}

var limitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
return current
`)

func (l *redisLimiter) Allow(key string, strict bool) bool {
	if l.client == nil || limitScript == nil {
		l.safeWarn("redis limiter disabled: client or script is nil")
		return !strict
	}

	result, err := l.allowByRedis(key)

	// 降级逻辑
	if err != nil {
		l.safeWarn("limiter fallback: redis error or circuit breaker open", slog.String("error", err.Error()))
		return !strict
	}

	return result.(bool)
}

func (l *redisLimiter) allowByRedis(key string) (any, error) {
	return l.cb.Execute(func() (any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), redisLimiterTimeout)
		defer cancel()

		count, err := limitScript.Run(ctx, l.client, []string{redisRateLimitKey(key)}, redisWindowMillis).Int64()
		if err != nil {
			return false, err
		}
		return count <= int64(l.burst), nil
	})
}

func (l *redisLimiter) safeWarn(msg string, attrs ...slog.Attr) {
	logger := l.logger
	if logger == nil {
		logger = slog.Default()
	}
	if len(attrs) == 0 {
		logger.Warn(msg)
		return
	}
	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}
	logger.Warn(msg, args...)
}

func redisRateLimitKey(key string) string {
	return "rate_limit:" + key
}

// HTTPRateLimitRedis 提供基于Redis的分布式限流
func HTTPRateLimitRedis(client redis.UniversalClient, b int, logger *slog.Logger) gin.HandlerFunc {
	store := NewRedisLimiter(client, b, logger)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !store.Allow(ip, false) {
			AbortWithError(c, errorHandler.ErrTooManyRequests)
			return
		}
		c.Next()
	}
}
