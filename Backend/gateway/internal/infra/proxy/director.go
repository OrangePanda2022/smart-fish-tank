package proxy

import (
	"gateway/internal/util"
	"net/http"
)

func NewDirector() func(*http.Request) {

	return func(req *http.Request) {
		req.Header.Set("X-Trace-ID", util.GenerateTraceID())
	}
}
