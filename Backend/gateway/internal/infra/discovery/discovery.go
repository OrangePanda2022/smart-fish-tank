package discovery

// TODO consul 动态服务发现

import (
	"log"
	"time"

	"gateway/internal/domain/upstream"
	"gateway/internal/infra/cache"
	"gateway/internal/util"

	"github.com/hashicorp/consul/api"
)

type ConsulDiscovery struct {
	client *api.Client
	cache  *cache.UpstreamCache
}

// NewConsulDiscovery 创建 Consul 客户端
func NewConsulDiscovery(consulAddr string) (*ConsulDiscovery, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulDiscovery{client: client}, nil
}

func (c *ConsulDiscovery) GetServiceUpstream(serviceName string) (*upstream.Upstream, error) {
	nodes, err := c.GetServiceNodes(serviceName)
	if err != nil {
		return nil, err
	}
	return util.ToCluster(serviceName, nodes)
}

// 获取服务健康节点
func (c *ConsulDiscovery) GetServiceNodes(serviceName string) ([]upstream.Node, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	nodes := make([]upstream.Node, 0, len(services))
	for _, s := range services {
		nodes = append(nodes, upstream.Node{
			Host: s.Service.Address,
			Port: s.Service.Port,
		})
	}
	return nodes, nil
}

// 获取服务健康节点
func (c *ConsulDiscovery) FetchServiceNodes(serviceName string) ([]upstream.Node, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	nodes := make([]upstream.Node, 0, len(services))
	for _, s := range services {
		nodes = append(nodes, upstream.Node{
			Host: s.Service.Address,
			Port: s.Service.Port,
		})
	}
	return nodes, nil
}

// WatchService 动态更新 cache
func (c *ConsulDiscovery) WatchService(serviceName string, interval time.Duration) {
	go func() {
		for {
			nodes, err := c.FetchServiceNodes(serviceName)
			if err != nil {
				log.Printf("Consul watch error: %v", err)
				time.Sleep(interval)
				continue
			}

			upstreamNodes := make([]upstream.Node, 0, len(nodes))
			for _, n := range nodes {
				upstreamNodes = append(upstreamNodes, upstream.Node{
					Host: n.Host,
					Port: n.Port,
				})
			}

			updateUpstream := upstream.Upstream{
				Name:  serviceName,
				Nodes: upstreamNodes,
			}

			c.cache.Set(serviceName, &updateUpstream)

			time.Sleep(interval)
		}
	}()
}
