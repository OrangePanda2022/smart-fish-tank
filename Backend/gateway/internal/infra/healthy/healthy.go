package healthy

import (
	"gateway/internal/domain/upstream"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Healthy struct {
	Client *http.Client
	Logger *zap.Logger
}

func NewHealthChecker(client *http.Client, logger *zap.Logger) *Healthy {
	return &Healthy{Client: client, Logger: logger}
}

// 过滤健康实例
func FilterHealthy(instances []upstream.Node) []*upstream.Node {
	var res []*upstream.Node
	for _, inst := range instances {
		if !inst.Healthy {
			res = append(res, &inst)
		}
	}
	return res
}

func (h *Healthy) IsHealthy(node upstream.Node) bool {
	// 初始化 sugar 日志
	sugar := h.Logger.Sugar()

	// 记录当前时间
	start := time.Now()

	resp, err := h.Client.Get(node.URL() + "/health")
	if err != nil {
		return false
	}

	latency := time.Since(start).Milliseconds()

	defer resp.Body.Close()

	sugar.Infow("health check", "node", node.URL(), "latency_ms", latency, "status", resp.StatusCode)

	return resp.StatusCode == http.StatusOK
}
