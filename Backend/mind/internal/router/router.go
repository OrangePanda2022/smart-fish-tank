package router

import (
	"context"

	"mind/internal/handler"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Register(h *server.Hertz, hdl *handler.Handler) {
	h.GET("/ping", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]string{"message": "pong"})
	})

	v1 := h.Group("/api/v1")
	v1.GET("/analyse/:tank_id", hdl.Analyse)

	// 预测 + MPC 端点
	v1.GET("/predict/:tank_id", hdl.Predict)
	v1.GET("/predict/:tank_id/stream", hdl.PredictStream)
	v1.GET("/predict/status", hdl.PredictStatus)
}
