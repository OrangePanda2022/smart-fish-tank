package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient 初始化 Redis 客户端，不可用时返回 nil 以触发降级策略
func NewRedisClient(ctx context.Context, addrs []string, password string, db int, logger *slog.Logger) redis.UniversalClient {
	if len(addrs) == 0 || (len(addrs) == 1 && addrs[0] == "") {
		logger.Warn("redis addresses empty, fallback to L1 cache")
		return nil
	}
	// UniversalOptions 会根据 Addrs 的数量自动选择 Client 类型
	opts := &redis.UniversalOptions{
		Addrs:        addrs,
		Password:     password,
		DB:           db, // Redis Cluster 默认只支持 DB 0
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	cli := redis.NewUniversalClient(opts)

	// 连接建立后立即 ping，确保启动阶段就能感知可用性并决定是否降级
	if err := cli.Ping(ctx).Err(); err != nil {
		logger.Warn("redis unavailable, fallback to L1 cache", slog.Any("addrs", addrs), slog.String("error", err.Error()))
		_ = cli.Close()
		return nil
	}
	logger.Info("redis connected", slog.Any("addrs", addrs))
	return cli
}
