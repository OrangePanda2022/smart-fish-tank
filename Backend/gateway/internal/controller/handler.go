package controller

import (
	"crypto/rand"
	"encoding/hex"
	"gateway/internal/aggregator"
	"gateway/internal/domain/request"
	"gateway/internal/infra/lb"
	"gateway/internal/infra/proxy"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GatewayHandler struct {
	router     *Router
	strict     bool
	proxy      *proxy.Proxy
	lb         *lb.StickyBalancer
	aggregator *aggregator.Aggregator
	logger     *zap.Logger
}

func NewGatewayHandler(router *Router, strict bool, proxy *proxy.Proxy, lb *lb.StickyBalancer, aggregator *aggregator.Aggregator, logger *zap.Logger) *GatewayHandler {
	return &GatewayHandler{
		router:     router,
		strict:     strict,
		proxy:      proxy,
		lb:         lb,
		aggregator: aggregator,
		logger:     logger,
	}
}

func (h *GatewayHandler) ProxyHandler(c *gin.Context) {
	sugar := h.logger.Sugar()
	sugar.Info("开始代理")
	// ?aggregate=true 判定是否是聚合请求
	if c.Query("aggregate") == "true" {
		if h.strict {
			// TODO 统一错误返回
			return
		}
		h.AggregateHandler(c)
		return
	}

	serviceName, exists := h.router.Match(c.Request.URL.Path)
	if exists != true {
		sugar.Info("未找到代理地址", c.Request.URL.Path)
		return
	}
	// 选出节点
	node, err := h.lb.Select(serviceName, RandomHex(4))
	if err != nil {
		sugar.Info("未找到节点", serviceName, "错误细节", err.Error())
		return
	}
	// 普通转发
	h.proxy.Forward(node, c.Writer, c.Request)
}

func RandomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *GatewayHandler) AggregateHandler(c *gin.Context) {
	services := strings.Split(c.Query("services"), ",")
	var reqs []request.AggregateRequest

	for _, s := range services {
		reqs = append(reqs, request.AggregateRequest{
			Key: s,
			// TODO 修正错误的 Path 拼接
			TargetURL: "http://" + s + "-service" + c.Request.URL.Path, // 自动拼接 Path
		})
	}

	// 调用聚合器
	results, err := h.aggregator.Aggregate(c.Request.Context(), reqs)
	if err != nil {
		// TODO 统一返回 统一错误处理
		c.JSON(500, gin.H{"error": "Aggregation failed", "details": err.Error()})
		return
	}

	c.JSON(200, results)
}

type Router struct {
	routes map[string]string
	logger *zap.Logger
}

func NewRouter(routes map[string]string, logger *zap.Logger) *Router {
	sugar := logger.Sugar()
	sugar.Info("初始化路由表")
	return &Router{routes: routes}
}

func (r *Router) Match(path string) (string, bool) {
	for prefix, serviceName := range r.routes {
		if strings.HasPrefix(path, prefix) {
			return serviceName, true
		}
	}
	return "", false
}
