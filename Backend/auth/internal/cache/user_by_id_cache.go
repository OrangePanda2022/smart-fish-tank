package cache

import (
	"auth/internal/domain/user"
	localCache "auth/internal/infra/cache"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

const (
	userByIDPrefix   = "user_by_id:"
	nullUserSentinel = "__NULL__"
)

type UserByIDMapper interface {
	Get(ctx context.Context, userID string) (u *user.UserEntity, hit bool, err error)
	Set(ctx context.Context, u *user.UserEntity) error
	SetNull(ctx context.Context, userID string)
	Delete(ctx context.Context, userID string)
}

// UserByIDCache 管理按 userID 的两级缓存，支持空值缓存防穿透
type UserByIDCache struct {
	l1      *localCache.LocalCache
	l2      redis.UniversalClient
	breaker *gobreaker.CircuitBreaker
	logger  *slog.Logger
	ttl     time.Duration
	nullTTL time.Duration
}

// NewUserByIDCache 初始化按 userID 的两级缓存组件
func NewUserByIDCache(l1 *localCache.LocalCache, l2 redis.UniversalClient, logger *slog.Logger, ttl, nullTTL time.Duration) *UserByIDCache {
	// 对 Redis 调用启用熔断，避免缓存层抖动放大到业务层
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-user-by-id",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("[ALERT] circuit state changed", slog.String("name", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	})
	return &UserByIDCache{l1: l1, l2: l2, breaker: cb, logger: logger, ttl: ttl, nullTTL: nullTTL}
}

// userByIDKey 统一 userID 缓存键格式
func userByIDCacheKey(userID string) string {
	return userByIDPrefix + userID
}

// Get 先查 L1，再查 L2，L2 命中后自动回填 L1
func (c *UserByIDCache) Get(ctx context.Context, userID string) (u *user.UserEntity, hit bool, err error) {
	key := userByIDCacheKey(userID)

	if payload, ok := c.l1.Get(key); ok {
		return c.decodePayload(payload)
	}

	if c.l2 == nil {
		return nil, false, nil
	}

	result, execErr := c.breaker.Execute(func() (interface{}, error) {
		return c.l2.Get(ctx, key).Result()
	})
	if execErr != nil {
		if execErr != redis.Nil {
			c.logger.Warn("redis get user by id cache failed", slog.String("error", execErr.Error()))
		}
		return nil, false, nil
	}

	payload, ok := result.(string)
	if !ok || payload == "" {
		return nil, false, nil
	}

	cachedUser, cacheHit, decodeErr := c.decodePayload(payload)
	if decodeErr != nil {
		c.logger.Warn("decode user cache payload failed", slog.String("error", decodeErr.Error()))
		return nil, false, nil
	}
	if !cacheHit {
		return nil, false, nil
	}

	writeTTL := c.ttl
	if payload == nullUserSentinel {
		writeTTL = c.nullTTL
	}
	if err := c.l1.Set(key, payload, writeTTL); err != nil {
		c.logger.Warn("write back l1 user cache failed", slog.String("error", err.Error()))
	}
	return cachedUser, true, nil
}

// Set 写入正常用户缓存到 L1/L2
func (c *UserByIDCache) Set(ctx context.Context, u *user.UserEntity) error {
	payloadBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}
	payload := string(payloadBytes)
	key := userByIDCacheKey(u.UserID)

	if err := c.l1.Set(key, payload, c.ttl); err != nil {
		return err
	}
	if c.l2 == nil {
		return nil
	}
	_, execErr := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Set(ctx, key, payload, c.ttl).Err()
	})
	if execErr != nil {
		c.logger.Warn("redis set user by id cache failed", slog.String("error", execErr.Error()))
	}
	return nil
}

// SetNull 写入空值缓存，防止缓存穿透
func (c *UserByIDCache) SetNull(ctx context.Context, userID string) {
	key := userByIDCacheKey(userID)
	if err := c.l1.Set(key, nullUserSentinel, c.nullTTL); err != nil {
		c.logger.Warn("set l1 null user cache failed", slog.String("error", err.Error()))
	}
	if c.l2 == nil {
		return
	}
	_, execErr := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Set(ctx, key, nullUserSentinel, c.nullTTL).Err()
	})
	if execErr != nil {
		c.logger.Warn("redis set null user cache failed", slog.String("error", execErr.Error()))
	}
}

// Delete 清理按 userID 的缓存
func (c *UserByIDCache) Delete(ctx context.Context, userID string) {
	key := userByIDCacheKey(userID)
	c.l1.Delete(key)
	if c.l2 == nil {
		return
	}
	_, execErr := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Del(ctx, key).Err()
	})
	if execErr != nil {
		c.logger.Warn("redis delete user cache failed", slog.String("error", execErr.Error()))
	}
}

func (c *UserByIDCache) decodePayload(payload string) (u *user.UserEntity, hit bool, err error) {
	// 空值哨兵代表明确不存在，用于阻断热点穿透
	if payload == "" {
		return nil, false, nil
	}
	if payload == nullUserSentinel {
		return nil, true, nil
	}
	var out user.UserEntity
	if err := json.Unmarshal([]byte(payload), &out); err != nil {
		return nil, false, err
	}
	return &out, true, nil
}
