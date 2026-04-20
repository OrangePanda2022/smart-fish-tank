package cache

import (
	"auth/internal/domain/user"
	errorHandler "auth/internal/error"
	localCache "auth/internal/infra/cache"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const revokedAtUserIDPrefix = "user_revoked:"

type WriteTask struct {
	UserID    string
	RevokedAt time.Time
}

// UserIDRevokedAtMapper 约束 UserID->RevokedAt 映射缓存能力
type RevokedAtUserIDMapper interface {
	Get(ctx context.Context, userID string, l1TTL time.Duration) (time.Time, bool, error)
	Set(ctx context.Context, userID string, revokedAt time.Time, ttl time.Duration) error
	Delete(ctx context.Context, userID string)
}

// UserIDRevokedAtCache 维护 UserID->RevokedAt 的两级缓存映射
type RevokedAtUserIDCache struct {
	l1            *localCache.LocalCache
	l2            redis.UniversalClient
	localDirtyMap sync.Map
	db            *gorm.DB
	breaker       *gobreaker.CircuitBreaker
	writeChan     chan WriteTask
	logger        *slog.Logger
}

func NewRevokedAtUserIDCache(l1 *localCache.LocalCache, l2 redis.UniversalClient, buffer int, db *gorm.DB, logger *slog.Logger) *RevokedAtUserIDCache {
	// Redis 不稳定时通过熔断快速失败，避免请求线程被慢调用拖垮
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-revoked-map",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("[ALERT] circuit state changed", slog.String("name", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	})
	return &RevokedAtUserIDCache{l1: l1, l2: l2, breaker: cb, writeChan: make(chan WriteTask, buffer), logger: logger}
}

// userIDKey 统一 Email 映射缓存键前缀，避免不同业务键冲突
func userIDKey(userID string) string {
	return revokedAtUserIDPrefix + userID
}

// Get 先查 L1，再查 L2，查到后回填 L1
func (c *RevokedAtUserIDCache) Get(ctx context.Context, userID string, l1TTL time.Duration) (time.Time, bool, error) {
	key := userIDKey(userID)
	if v, ok := c.l1.Get(key); ok {
		revoked, _ := strconv.ParseInt(v, 10, 64)
		return time.UnixMicro(revoked), true, nil
	}
	if c.l2 == nil {
		return time.Time{}, false, nil
	}
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.l2.Get(ctx, key).Result()
	})
	if err != nil {
		if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
			return time.Time{}, false, errorHandler.ErrCircuitOpen.WithErr(err)
		}
		if err == redis.Nil {
			return time.Time{}, false, nil
		}
		c.logger.Warn("redis get revoked map failed", slog.String("error", err.Error()))
		return time.Time{}, false, nil
	}
	value, _ := result.(string)
	if value == "" {
		return time.Time{}, false, nil
	}
	_ = c.l1.Set(key, value, l1TTL)

	data, _ := strconv.ParseInt(value, 10, 64)
	return time.UnixMicro(data), true, nil
}

// Set 同时写入 L1/L2，L2 失败不会影响主流程
func (c *RevokedAtUserIDCache) Set(ctx context.Context, userID string, revokedAt time.Time, ttl time.Duration) error {
	key := userIDKey(userID)
	revoked := strconv.FormatInt(revokedAt.UnixMicro(), 10)
	if err := c.l1.Set(key, revoked, ttl); err != nil {
		return fmt.Errorf("set l1 email map failed: %w", err)
	}
	if c.l2 == nil {
		return nil
	}
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.l2.Set(ctx, key, revokedAt, ttl).Err()
	})
	if err != nil {
		c.logger.Warn("redis set email map failed", slog.String("error", err.Error()))
	}

	select {
	case c.writeChan <- WriteTask{UserID: userID, RevokedAt: revokedAt}:
	default:
		c.logger.Warn("write chan full, skip db update", slog.String("userID", userID))
	}
	return nil
}

// Delete 清理映射缓存，防止脏读
func (c *RevokedAtUserIDCache) Delete(ctx context.Context, userID string) {
	key := userIDKey(userID)
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

func (c *RevokedAtUserIDCache) flushWorker() {
	ticker := time.NewTicker(5 * time.Second)
	pending := make(map[string]time.Time)

	for {
		select {
		case task := <-c.writeChan:
			pending[task.UserID] = task.RevokedAt
		case <-ticker.C:
			if len(pending) == 0 {
				continue
			}

			if err := c.batchSyncToDB(pending); err != nil {
				c.logger.Error("batch sync to db failed", slog.String("error", err.Error()))
			}
			pending = make(map[string]time.Time) // 清空
		}
	}
}

func (c *RevokedAtUserIDCache) batchSyncToDB(updates map[string]time.Time) error {
	var users []user.UserEntity
	for id, revokedAt := range updates {
		users = append(users, user.UserEntity{UserID: id, UserRevoked: revokedAt})
	}

	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_revoked"}),
	}).Create(&users).Error
}
