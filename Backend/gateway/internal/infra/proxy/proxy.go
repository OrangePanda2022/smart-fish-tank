package proxy

import (
	"context"
	"encoding/json"
	"gateway/internal/domain/upstream"
	"gateway/internal/infra/buffer"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

type Proxy struct {
	proxy  *httputil.ReverseProxy
	logger *zap.Logger
}

func NewProxy(director func(*http.Request), buffer *buffer.BytePool, logger *zap.Logger) *Proxy {

	proxy := &httputil.ReverseProxy{
		Director:   director,
		BufferPool: buffer,
		Transport:  transport,
	}

	return &Proxy{
		proxy:  proxy,
		logger: logger,
	}
}

// 把请求转发到 upstraem 上游服务
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

// 选择一个上游节点，并且转发服务到那里
func (p *Proxy) Forward(node *upstream.Node, w http.ResponseWriter, r *http.Request) {
	p.logger.Sugar().Info("收到请求", node.URL())

	r.Header.Set("X-Forwarded-For", r.RemoteAddr)

	targetURL, _ := url.Parse(node.URL())

	// 手动修改请求目标，而不重新创建 Proxy
	r.URL.Scheme = targetURL.Scheme
	r.URL.Host = targetURL.Host
	// 后端服务校验 Host 时必须匹配
	r.Host = targetURL.Host

	// 执行转发
	p.proxy.ServeHTTP(w, r)
}

func (p *Proxy) DoRequest(ctx context.Context, targetURL string) (interface{}, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", targetURL, nil)

	// 复用之前的 Transport
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取数据并解析为 JSON
	body, _ := io.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return data, nil
}
