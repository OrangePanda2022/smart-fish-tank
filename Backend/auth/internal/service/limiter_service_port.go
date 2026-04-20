package service

// LimiterService 定义按 key 判断是否放行的限流能力
type LimiterService interface {
	Allow(key string) bool
}
