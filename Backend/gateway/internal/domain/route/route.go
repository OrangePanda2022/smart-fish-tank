package route

import "gateway/internal/domain/upstream"

type RouteID string

type Route struct {
	ID       RouteID
	Name     string
	Path     string
	Method   string
	Upstream upstream.UpstreamID
	Plugins  []string
	Enabled  bool
}

func NewRoute(
	id RouteID,
	name string,
	path string,
	method string,
	upstreamID upstream.UpstreamID,
) *Route {
	return &Route{
		ID:       id,
		Name:     name,
		Path:     path,
		Method:   method,
		Upstream: upstreamID,
		Enabled:  true,
	}
}

func (r *Route) Disable() {
	r.Enabled = false
}

func (r *Route) Enable() {
	r.Enabled = true
}
