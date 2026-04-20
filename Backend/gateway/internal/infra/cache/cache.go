package cache

import (
	"encoding/json"
	"gateway/internal/domain/upstream"

	"github.com/VictoriaMetrics/fastcache"
)

type UpstreamCache struct {
	fc *fastcache.Cache
}

func NewUpstreamCache(maxBytes int) *UpstreamCache {
	return &UpstreamCache{
		fc: fastcache.New(maxBytes),
	}
}

// Set 将 Upstream 存入缓存
func (c *UpstreamCache) Set(serviceName string, up *upstream.Upstream) error {
	data, err := json.Marshal(up)
	if err != nil {
		return err
	}
	c.fc.Set([]byte(serviceName), data)
	return nil
}

// Get 从缓存获取 Upstream
func (c *UpstreamCache) Get(serviceName string) (*upstream.Upstream, bool) {
	data := c.fc.Get(nil, []byte(serviceName))
	if data == nil {
		return nil, false
	}

	var up upstream.Upstream
	if err := json.Unmarshal(data, &up); err != nil {
		return nil, false
	}
	return &up, true
}

// Delete 删除缓存
func (c *UpstreamCache) Delete(serviceName string) {
	c.fc.Del([]byte(serviceName))
}
