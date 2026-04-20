package handler

import (
	"net/http"
	"strconv"

	"sensor/internal/model"
	"sensor/internal/service"

	"github.com/gin-gonic/gin"
)

// SensorHandler 表示传感器的HTTP处理器
type SensorHandler struct {
	service *service.SensorService // 传感器服务
}

// NewSensorHandler 创建传感器处理器
func NewSensorHandler(svc *service.SensorService) *SensorHandler {
	return &SensorHandler{service: svc}
}

// GetLatestByTank 处理 GET /api/tanks/:tankId/sensors/latest
func (h *SensorHandler) GetLatestByTank(c *gin.Context) {
	tankID := c.Param("tankId")
	if tankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tankId不能为空"})
		return
	}

	data, err := h.service.GetLatestByTank(tankID)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到数据"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO
	dto := model.ToDTO(
		data.DeviceID,
		data.TankID,
		data.Temperature,
		data.PH,
		data.Oxygen,
		data.Ammonia,
		data.WaterLevel,
		data.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)

	c.JSON(http.StatusOK, dto)
}

// GetHistoryByTank 处理 GET /api/tanks/:tankId/sensors/history
func (h *SensorHandler) GetHistoryByTank(c *gin.Context) {
	tankID := c.Param("tankId")
	if tankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tankId不能为空"})
		return
	}

	// 获取查询参数
	start := c.Query("start")
	end := c.Query("end")
	limitStr := c.Query("limit")

	limit := 100
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	data, err := h.service.GetHistoryByTank(tankID, start, end, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表
	dtos := make([]model.SensorDataDTO, len(data))
	for i, d := range data {
		dtos[i] = model.ToDTO(
			d.DeviceID,
			d.TankID,
			d.Temperature,
			d.PH,
			d.Oxygen,
			d.Ammonia,
			d.WaterLevel,
			d.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		)
	}

	result := model.SensorDataListDTO{
		Data:  dtos,
		Count: len(dtos),
	}

	c.JSON(http.StatusOK, result)
}

// GetByDevice 处理 GET /api/devices/:deviceId/sensors
func (h *SensorHandler) GetByDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deviceId不能为空"})
		return
	}

	// 获取查询参数
	limitStr := c.Query("limit")

	limit := 100
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	data, err := h.service.GetByDevice(deviceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为DTO列表
	dtos := make([]model.SensorDataDTO, len(data))
	for i, d := range data {
		dtos[i] = model.ToDTO(
			d.DeviceID,
			d.TankID,
			d.Temperature,
			d.PH,
			d.Oxygen,
			d.Ammonia,
			d.WaterLevel,
			d.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		)
	}

	result := model.SensorDataListDTO{
		Data:  dtos,
		Count: len(dtos),
	}

	c.JSON(http.StatusOK, result)
}

// HealthCheck 处理 GET /health
func (h *SensorHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "sensor-service",
	})
}
