package util

import "gateway/internal/domain/upstream"

func ToCluster(serviceName string, nodes []upstream.Node) (*upstream.Upstream, error) {
	return &upstream.Upstream{
		Name:  serviceName,
		Nodes: nodes,
	}, nil
}
