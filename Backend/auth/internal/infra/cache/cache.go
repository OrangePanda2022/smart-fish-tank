package localCache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
)

// localCache 使用 BigCache 作为 L1 本地缓存，并自行维护过期时间
type LocalCache struct {
	bc *bigcache.BigCache
}

// NewLocalCache 初始化 BigCache 实例，作为统一 L1 缓存基础组件
func NewLocalCache(defaultLifeWindow time.Duration) (*LocalCache, error) {
	// BigCache 默认以 LifeWindow 做分片回收；实际 TTL 由 value 中的过期时间控制
	_ = defaultLifeWindow
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("new bigcache failed: %w", err)
	}
	return &LocalCache{bc: cache}, nil
}

// Set 将过期时间与 value 拼接存储，避免依赖 BigCache 全局 TTL
func (l *LocalCache) Set(key, value string, ttl time.Duration) error {
	expireAt := time.Now().Add(ttl).Unix()
	payload := strconv.FormatInt(expireAt, 10) + "|" + value
	return l.bc.Set(key, []byte(payload))
}

// Get 读取并校验自定义过期时间，过期数据会被懒删除
func (l *LocalCache) Get(key string) (string, bool) {
	data, err := l.bc.Get(key)
	if err != nil {
		return "", false
	}
	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return "", false
	}
	expireAt, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return "", false
	}
	if time.Now().Unix() > expireAt {
		_ = l.bc.Delete(key)
		return "", false
	}
	return parts[1], true
}

// Delete 删除指定 key（忽略不存在错误）
func (l *LocalCache) Delete(key string) {
	_ = l.bc.Delete(key)
}

// Close 释放 BigCache 占用资源
func (l *LocalCache) Close() error {
	return l.bc.Close()
}
