package lb

import (
	"fmt"
	"gateway/internal/domain/upstream"
	"gateway/internal/infra/cache"
	"gateway/internal/infra/discovery"
	"gateway/internal/infra/healthy"
	"gateway/internal/util"

	"go.uber.org/zap"
)

type StickyBalancer struct {
	discovery *discovery.ConsulDiscovery
	cache     *cache.UpstreamCache
	logger    *zap.Logger
}

func NewSticktBalancer(discovery *discovery.ConsulDiscovery, cache *cache.UpstreamCache, logger *zap.Logger) *StickyBalancer {
	return &StickyBalancer{
		discovery: discovery,
		cache:     cache,
		logger:    logger,
	}
}

func (b *StickyBalancer) Select(serviceName string, userKey string) (*upstream.Node, error) {
	// 从缓存获取数据
	cluster, exists := b.cache.Get(serviceName)
	if exists == false {
		upstream, err := b.discovery.GetServiceUpstream(serviceName)
		if err != nil {
			return nil, err
		}
		cluster = upstream
	}
	// 获取所有健康实例
	healthy := healthy.FilterHealthy(cluster.Nodes)
	if len(healthy) == 0 {
		return nil, fmt.Errorf("ErrNoInstanceAvailable")
	}

	// 使用原生 JumpHash 计算索引
	index := util.JumpHash(userKey, len(healthy))

	return healthy[index], nil
}
