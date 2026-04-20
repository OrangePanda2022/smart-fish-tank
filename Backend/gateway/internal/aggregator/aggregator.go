package aggregator

import (
	"context"
	"gateway/internal/domain/request"
	"gateway/internal/infra/proxy"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Aggregator struct {
	proxy *proxy.Proxy
}

func NewAggregator(proxy *proxy.Proxy) *Aggregator {
	return &Aggregator{proxy: proxy}
}

// Aggregate 执行聚合请求
func (a *Aggregator) Aggregate(ctx context.Context, requests []request.AggregateRequest) (map[string]interface{}, error) {
	group, ctx := errgroup.WithContext(ctx)
	results := make(map[string]interface{})
	var mu sync.Mutex

	for _, req := range requests {
		group.Go(func() error {
			// 调用 Infra 层的 Proxy 逻辑
			resp, err := a.proxy.DoRequest(ctx, req.TargetURL)
			if err != nil {
				return err
			}

			// 合并结果
			mu.Lock()
			results[req.Key] = resp
			mu.Unlock()
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}
