package proxy

import (
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{
	// 总最大空闲连接数
	MaxIdleConns: 10000,
	// 每个最大空闲连接
	MaxIdleConnsPerHost: 500,
	// 空闲连接多久关闭
	IdleConnTimeout: 90 * time.Second,

	// 超时控制
	DialContext: (&net.Dialer{
		// 建立 TCP 连接的超时
		Timeout: 5 * time.Second,
		// 探测存活频率
		KeepAlive: 30 * time.Second,
	}).DialContext,

	// HTTPS 握手超时
	TLSHandshakeTimeout: 5 * time.Second,
	// 等待后端 Header 返回的超时（mind ReAct 多轮调工具耗时长，留足余量）
	ResponseHeaderTimeout: 120 * time.Second,

	// 允许 Gzip，减少带宽
	DisableCompression: false,
	// 使用 HTTP/2
	ForceAttemptHTTP2: true,
}
