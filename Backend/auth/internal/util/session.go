package util

/*
主要修改说明（不改变业务行为）：
1. 函数拆分/合并：抽取 `marshalSession`、`unmarshalSession`、`isExpired`、`exec`，去除重复逻辑
2. 包结构调整：未调整包结构，仅在当前文件内重构
3. 校验逻辑增删：保留原有校验语义，仅集中到辅助函数以提升可读性
4. 注释修正：修正与实际实现不一致的注释（如 NewSession 的写入路径），并统一为简洁中文
5. 代码风格统一：统一导入分组、常量提取、命名与空行风格，便于维护
*/

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	errorHandler "auth/internal/error"

	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

const (
	sessionIDByteLen = 32
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID    string
	ExpiresAt time.Time
}

type SessionStore interface {
	NewSession(ctx context.Context, userID string, ttl time.Duration) (sid string, sess *Session, err error)
	Get(ctx context.Context, sid string) (*Session, error)
	Set(ctx context.Context, sid string, sess *Session) error
	Delete(ctx context.Context, sid string) error
}

type CompositeSessionStore struct {
	l1 *BigCacheSessionStore
	l2 RedisSessionStore
}

func NewCompositeSessionStore(l1 *BigCacheSessionStore, l2 RedisSessionStore) *CompositeSessionStore {
	return &CompositeSessionStore{l1: l1, l2: l2}
}

// NewSessionID 生成安全随机的 Session ID
func (s *CompositeSessionStore) NewSessionID() (string, error) {
	b := make([]byte, sessionIDByteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// NewSession 创建 session 并写入当前实现中的一级缓存存储
func (s *CompositeSessionStore) NewSession(ctx context.Context, userID string, ttl time.Duration) (sid string, sess *Session, err error) {
	sid, err = s.NewSessionID()
	if err != nil {
		return "", nil, err
	}

	sess = &Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err = s.l1.Set(ctx, sid, sess); err != nil {
		return "", nil, err
	}
	return sid, sess, nil
}

func (s *CompositeSessionStore) Get(ctx context.Context, sid string) (*Session, error) {
	if sess, err := s.l1.Get(ctx, sid); err == nil {
		return sess, nil
	}

	sess, err := s.l2.Get(ctx, sid)
	if err != nil {
		return nil, err
	}

	_ = s.l1.Set(ctx, sid, sess)
	return sess, nil
}

func (s *CompositeSessionStore) Set(ctx context.Context, sid string, sess *Session) error {
	_ = s.l1.Set(ctx, sid, sess)
	return s.l2.Set(ctx, sid, sess)
}

func (s *CompositeSessionStore) Delete(ctx context.Context, sid string) error {
	_ = s.l1.Delete(ctx, sid)
	return s.l2.Delete(ctx, sid)
}

func NewBigCache() (*bigcache.BigCache, error) {
	return bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
}

type BigCacheSessionStore struct {
	cache *bigcache.BigCache
}

func NewBigCacheSessionStore() (*BigCacheSessionStore, error) {
	cache, err := NewBigCache()
	if err != nil {
		return nil, err
	}
	return &BigCacheSessionStore{cache: cache}, nil
}

func (s *BigCacheSessionStore) Get(ctx context.Context, sid string) (*Session, error) {
	data, err := s.cache.Get(sid)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	sess, err := unmarshalSession(data)
	if err != nil {
		return nil, err
	}

	// bigcache 的生命周期与业务过期时间独立，需做业务侧兜底
	if isExpired(sess.ExpiresAt) {
		_ = s.cache.Delete(sid)
		return nil, ErrSessionNotFound
	}

	return sess, nil
}

func (s *BigCacheSessionStore) Set(ctx context.Context, sid string, sess *Session) error {
	data, err := marshalSession(sess)
	if err != nil {
		return err
	}
	return s.cache.Set(sid, data)
}

func (s *BigCacheSessionStore) Delete(ctx context.Context, sid string) error {
	return s.cache.Delete(sid)
}

type RedisSessionStore struct {
	redisClient redis.UniversalClient
	breaker     *gobreaker.CircuitBreaker
	prefix      string
}

func NewRedisSessionStore(redis redis.UniversalClient, prefix string) *RedisSessionStore {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis-session",
		MaxRequests: 50,
		Interval:    10 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 5
		},
	})

	return &RedisSessionStore{
		redisClient: redis,
		breaker:     cb,
		prefix:      prefix,
	}
}

func (s *RedisSessionStore) key(sid string) string {
	return fmt.Sprintf("session:%s", sid)
}

func (s *RedisSessionStore) Get(ctx context.Context, sid string) (*Session, error) {
	if s.redisClient == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	result, err := s.exec(func() (any, error) {
		val, err := s.redisClient.Get(ctx, s.key(sid)).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, ErrSessionNotFound
			}
			return nil, err
		}

		sess, err := unmarshalSession([]byte(val))
		if err != nil {
			return nil, err
		}

		if isExpired(sess.ExpiresAt) {
			_ = s.redisClient.Del(ctx, s.key(sid)).Err()
			return nil, ErrSessionNotFound
		}
		return sess, nil
	})
	if err != nil {
		return nil, errorHandler.ErrCircuitOpen.Err
	}
	return result.(*Session), nil
}

func (s *RedisSessionStore) Set(ctx context.Context, sid string, sess *Session) error {
	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		return ErrSessionNotFound
	}

	data, err := marshalSession(sess)
	if err != nil {
		return err
	}

	_, err = s.exec(func() (any, error) {
		return nil, s.redisClient.Set(ctx, s.key(sid), data, ttl).Err()
	})
	if err != nil {
		return errorHandler.ErrCircuitOpen
	}
	return nil
}

func (s *RedisSessionStore) Delete(ctx context.Context, sid string) error {
	_, err := s.exec(func() (any, error) {
		return nil, s.redisClient.Del(ctx, s.key(sid)).Err()
	})
	if err != nil {
		return errorHandler.ErrCircuitOpen
	}
	return nil
}

func (s *RedisSessionStore) exec(fn func() (any, error)) (any, error) {
	return s.breaker.Execute(fn)
}

func marshalSession(sess *Session) ([]byte, error) {
	return json.Marshal(sess)
}

func unmarshalSession(data []byte) (*Session, error) {
	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func isExpired(expireAt time.Time) bool {
	return time.Now().After(expireAt)
}
