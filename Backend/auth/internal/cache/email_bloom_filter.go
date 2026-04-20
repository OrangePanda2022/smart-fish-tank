package cache

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

// EmailBloomFilter 提供邮箱存在性快速判定能力
type EmailBloomFilter interface {
	MightContain(ctx context.Context, email string) (bool, error)
	Add(ctx context.Context, email string) error
}

// RedisEmailBloomFilter 基于 RedisBloom 模块实现邮箱布隆过滤器
type RedisEmailBloomFilter struct {
	client  redis.UniversalClient
	key     string
	breaker *gobreaker.CircuitBreaker
	logger  *slog.Logger
}

// NewEmailBloomFilter 创建邮箱布隆过滤器；Redis 不可用时自动降级为 no-op
func NewEmailBloomFilter(client redis.UniversalClient, logger *slog.Logger, key string, errorRate float64, capacity uint64) EmailBloomFilter {
	if client == nil {
		logger.Warn("redis bloom disabled, fallback to noop")
		return &noopEmailBloomFilter{}
	}
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-email-bloom",
		MaxRequests: 3,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warn("[ALERT] circuit state changed", slog.String("name", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	})
	f := &RedisEmailBloomFilter{client: client, key: key, breaker: cb, logger: logger}
	f.reserve(context.Background(), errorRate, capacity)
	return f
}

func (f *RedisEmailBloomFilter) reserve(ctx context.Context, errorRate float64, capacity uint64) {
	_, err := f.client.Do(ctx, "BF.RESERVE", f.key, errorRate, capacity).Result()
	if err == nil {
		f.logger.Info("redis bloom reserved", slog.String("key", f.key), slog.Float64("error_rate", errorRate), slog.Uint64("capacity", capacity))
		return
	}
	if strings.Contains(strings.ToLower(err.Error()), "exists") {
		return
	}
	f.logger.Warn("reserve redis bloom failed, fallback to best effort", slog.String("error", err.Error()), slog.String("key", f.key))
}

// MightContain 返回布隆过滤器是否可能包含该邮箱
func (f *RedisEmailBloomFilter) MightContain(ctx context.Context, email string) (bool, error) {
	result, err := f.breaker.Execute(func() (interface{}, error) {
		return f.client.Do(ctx, "BF.EXISTS", f.key, email).Result()
	})
	if err != nil {
		return false, err
	}
	return parseRedisBool(result)
}

// Add 将邮箱写入布隆过滤器
func (f *RedisEmailBloomFilter) Add(ctx context.Context, email string) error {
	_, err := f.breaker.Execute(func() (interface{}, error) {
		return f.client.Do(ctx, "BF.ADD", f.key, email).Result()
	})
	return err
}

// noopEmailBloomFilter 作为降级实现，不阻断主业务
type noopEmailBloomFilter struct{}

func (n *noopEmailBloomFilter) MightContain(_ context.Context, _ string) (bool, error) {
	return true, nil
}

func (n *noopEmailBloomFilter) Add(_ context.Context, _ string) error {
	return nil
}

func parseRedisBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case int64:
		return val == 1, nil
	case bool:
		return val, nil
	case string:
		return val == "1" || strings.EqualFold(val, "true"), nil
	default:
		return false, fmt.Errorf("unexpected redis bool type: %T", v)
	}
}
