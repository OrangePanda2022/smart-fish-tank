package main

import (
	"fmt"
	"gateway/internal/aggregator"
	"gateway/internal/config"
	"gateway/internal/controller"
	"gateway/internal/domain/upstream"
	"gateway/internal/infra/buffer"
	"gateway/internal/infra/cache"
	"gateway/internal/infra/discovery"
	"gateway/internal/infra/lb"
	"gateway/internal/infra/proxy"
	"gateway/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	cache := cache.NewUpstreamCache(1024)
	logger, err := logger.InitLogger()
	if err != nil {
	}
	Init(cache, "./server-map.yaml", logger)
	routeMap, err := config.LoadRouteMap("./route-map.yaml")
	if err != nil {
		fmt.Println("加载路由配置失败:", err.Error())
		return
	}
	router := controller.NewRouter(routeMap, logger)
	director := proxy.NewDirector()
	buffer := buffer.NewBufferPool()
	proxy := proxy.NewProxy(director, buffer, logger)
	discovery, err := discovery.NewConsulDiscovery("")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lb := lb.NewSticktBalancer(discovery, cache, logger)
	aggregator := aggregator.NewAggregator(proxy)
	handler := controller.NewGatewayHandler(router, false, proxy, lb, aggregator, logger)
	r.Any("/*action", handler.ProxyHandler)

	r.Run(":8080")
}

func Init(cache *cache.UpstreamCache, fallbackPath string, logger *zap.Logger) {
	// 加载兜底配置
	fallback, err := config.LoadFallback(fallbackPath, logger)
	if err == nil {
		// 将兜底的 Upstream 直接写入 FastCache
		for serviceName, nodes := range fallback.Upstreams {
			cache.Set(serviceName, &upstream.Upstream{
				Name:  serviceName,
				Nodes: nodes,
			})
		}
		logger.Info("已加载静态兜底路由与上游配置")
	} else {
		fmt.Println(err.Error())
	}
}
