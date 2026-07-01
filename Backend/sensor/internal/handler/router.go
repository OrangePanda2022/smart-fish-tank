package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册路由
func (h *SensorHandler) RegisterRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", h.HealthCheck)

	// API路由
	api := router.Group("/api/v1")
	{
		// 鱼缸传感器路由
		tanks := api.Group("/tanks/:tankId")
		{
			tanks.GET("/sensors/latest", h.GetLatestByTank)
			tanks.GET("/sensors/history", h.GetHistoryByTank)
			tanks.POST("/sensors", h.CreateSensorData)
		}

		// 设备传感器路由
		devices := api.Group("/devices/:deviceId")
		{
			devices.GET("/sensors", h.GetByDevice)
		}
	}
}
