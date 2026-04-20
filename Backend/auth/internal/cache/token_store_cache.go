package cache

import (
	localCache "auth/internal/infra/cache"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

const (
	blackJTIPrefix = "black_jti:"
	resetPrefix    = "reset_token:"
)

// TokenCache 约束 Refresh 黑名单与 resetToken 白名单能力
type TokenCache interface {
	blacklistJTI(ctx context.Context, jti string, ttl time.Duration)
	isBlacklisted(ctx context.Context, jti string) bool
	PutResetToken(ctx context.Context, token, userID string, ttl time.Duration)
	ConsumeResetToken(ctx context.Context, token string) (string, bool)
}

// TokenStore 管理 Refresh JTI 黑名单与 resetToken 白名单
type TokenStore struct {
	l1      *localCache.LocalCache
	l2      redis.UniversalClient
	breaker *gobreaker.CircuitBreaker
	logger  *slog.Logger
}

func NewTokenStore(l1 *localCache.LocalCache, l2 redis.UniversalClient, logger *slog.Logger) *TokenStore {
	// Token 黑名单与 resetToken 都属于认证核心路径，熔断可避免 Redis 级联故障
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-token-store",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("[ALERT] circuit state changed", slog.String("name", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	})
	return &TokenStore{l1: l1, l2: l2, breaker: cb, logger: logger}
}

// BlacklistJTI 将 Refresh Token 的 JTI 拉黑，用于登出与敏感操作后的失效控制
func (s *TokenStore) blacklistJTI(ctx context.Context, jti string, ttl time.Duration) {
	key := blackJTIPrefix + jti
	_ = s.l1.Set(key, "1", ttl)
	if s.l2 == nil {
		return
	}
	_, err := s.breaker.Execute(func() (interface{}, error) {
		return nil, s.l2.Set(ctx, key, "1", ttl).Err()
	})
	if err != nil {
		s.logger.Warn("redis blacklist jti failed", slog.String("error", err.Error()))
	}
}

// IsBlacklisted 判断 JTI 是否已失效
func (s *TokenStore) isBlacklisted(ctx context.Context, jti string) bool {
	key := blackJTIPrefix + jti
	if _, ok := s.l1.Get(key); ok {
		return true
	}
	if s.l2 == nil {
		return false
	}
	result, err := s.breaker.Execute(func() (interface{}, error) {
		return s.l2.Get(ctx, key).Result()
	})
	if err != nil {
		if err != redis.Nil {
			s.logger.Warn("redis check black jti failed", slog.String("error", err.Error()))
		}
		return false
	}
	if str, ok := result.(string); ok && str != "" {
		_ = s.l1.Set(key, "1", 10*time.Minute)
		return true
	}
	return false
}

// PutResetToken 将重置密码令牌放入白名单缓存
func (s *TokenStore) PutResetToken(ctx context.Context, token, userID string, ttl time.Duration) {
	key := resetPrefix + token
	_ = s.l1.Set(key, userID, ttl)
	if s.l2 == nil {
		return
	}
	_, err := s.breaker.Execute(func() (interface{}, error) {
		return nil, s.l2.Set(ctx, key, userID, ttl).Err()
	})
	if err != nil {
		s.logger.Warn("redis set reset token failed", slog.String("error", err.Error()))
	}
}

// ConsumeResetToken 消费一次性 resetToken，读取成功后立即删除，避免重放
func (s *TokenStore) ConsumeResetToken(ctx context.Context, token string) (string, bool) {
	key := resetPrefix + token
	if userID, ok := s.l1.Get(key); ok {
		s.l1.Delete(key)
		if s.l2 != nil {
			_, _ = s.breaker.Execute(func() (interface{}, error) {
				return nil, s.l2.Del(ctx, key).Err()
			})
		}
		return userID, true
	}
	if s.l2 == nil {
		return "", false
	}
	result, err := s.breaker.Execute(func() (interface{}, error) {
		return s.l2.GetDel(ctx, key).Result()
	})
	if err != nil {
		if err != redis.Nil {
			s.logger.Warn("redis consume reset token failed", slog.String("error", err.Error()))
		}
		return "", false
	}
	userID, _ := result.(string)
	if userID == "" {
		return "", false
	}
	s.l1.Delete(key)
	return userID, true
}
