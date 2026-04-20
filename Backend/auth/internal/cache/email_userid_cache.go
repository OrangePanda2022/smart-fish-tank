package cache

import (
	errorHandler "auth/internal/error"
	localCache "auth/internal/infra/cache"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

const emailUserIDPrefix = "email_uid:"

// EmailUserIDMapper 约束 Email->UserID 映射缓存能力
type EmailUserIDMapper interface {
	Get(ctx context.Context, email string, l1TTL time.Duration) (string, bool, error)
	Set(ctx context.Context, email, userID string, ttl time.Duration) error
	Delete(ctx context.Context, email string)
}

// EmailUserIDCache 维护 Email -> UserID 的两级缓存映射
type EmailUserIDCache struct {
	l1      *localCache.LocalCache
	l2      redis.UniversalClient
	breaker *gobreaker.CircuitBreaker
	logger  *slog.Logger
}

func NewEmailUserIDCache(l1 *localCache.LocalCache, l2 redis.UniversalClient, logger *slog.Logger) *EmailUserIDCache {
	// Redis 不稳定时通过熔断快速失败，避免请求线程被慢调用拖垮
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-email-map",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("[ALERT] circuit state changed", slog.String("name", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	})
	return &EmailUserIDCache{l1: l1, l2: l2, breaker: cb, logger: logger}
}

// emailKey 统一 Email 映射缓存键前缀，避免不同业务键冲突
func emailKey(email string) string {
	return emailUserIDPrefix + email
}

// Get 先查 L1，再查 L2，查到后回填 L1
func (c *EmailUserIDCache) Get(ctx context.Context, email string, l1TTL time.Duration) (string, bool, error) {
	key := emailKey(email)
	if v, ok := c.l1.Get(key); ok {
		return v, true, nil
	}
	if c.l2 == nil {
		return "", false, nil
	}
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.l2.Get(ctx, key).Result()
	})
	if err != nil {
		if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
			return "", false, errorHandler.ErrCircuitOpen.WithErr(err)
		}
		if err == redis.Nil {
			return "", false, nil
		}
		c.logger.Warn("redis get email map failed", slog.String("error", err.Error()))
		return "", false, nil
	}
	value, _ := result.(string)
	if value == "" {
		return "", false, nil
	}
	_ = c.l1.Set(key, value, l1TTL)
	return value, true, nil
}

// Set 同时写入 L1/L2，L2 失败不会影响主流程
func (c *EmailUserIDCache) Set(ctx context.Context, email, userID string, ttl time.Duration) error {
	key := emailKey(email)
	if err := c.l1.Set(key, userID, ttl); err != nil {
		return fmt.Errorf("set l1 email map failed: %w", err)
	}
	if c.l2 == nil {
		return nil
	}
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Set(ctx, key, userID, ttl).Err()
	})
	if err != nil {
		c.logger.Warn("redis set email map failed", slog.String("error", err.Error()))
	}
	return nil
}

// Delete 清理映射缓存，防止脏读
func (c *EmailUserIDCache) Delete(ctx context.Context, email string) {
	key := emailKey(email)
	c.l1.Delete(key)
	if c.l2 == nil {
		return
	}
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Del(ctx, key).Err()
	})
	if err != nil {
		c.logger.Warn("redis delete email map failed", slog.String("error", err.Error()))
	}
}
