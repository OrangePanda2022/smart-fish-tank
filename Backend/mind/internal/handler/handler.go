package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"mind/internal/domain"
	"mind/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Handler struct {
	analyseSvc *service.AnalyseService
	predictSvc *service.PredictService
}

func NewHandler(analyseSvc *service.AnalyseService, predictSvc *service.PredictService) *Handler {
	return &Handler{
		analyseSvc: analyseSvc,
		predictSvc: predictSvc,
	}
}

func (h *Handler) Analyse(_ context.Context, c *app.RequestContext) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tank_id is required",
		})
		return
	}

	var req domain.SensorDataListDTO
	if err := c.Bind(&req); err != nil {
		// 允许空请求，继续使用数据库数据
	}

	var sensorData []*domain.SensorData
	if len(req.Data) > 0 {
		sensorData = convertToSensorData(req.Data)
	}

	// 判断流式请求
	if c.Query("stream") == "true" {
		h.analyseStream(c, tankID, sensorData)
		return
	}

	response, err := h.analyseSvc.Analyse(context.Background(), tankID, sensorData)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": fmt.Sprintf("analyse failed: %v", err),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// Predict 处理 GET /api/v1/predict/:tank_id
// 调用方只需传入 tank_id，可选 horizon 查询参数
func (h *Handler) Predict(_ context.Context, c *app.RequestContext) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tank_id is required",
		})
		return
	}

	// 默认预测12步（1小时 @ 5分钟步长），最多60步
	horizon := 12
	if v, err := strconv.Atoi(c.Query("horizon")); err == nil && v > 0 && v <= 60 {
		horizon = v
	}

	resp, err := h.predictSvc.Predict(context.Background(), tankID, horizon)
	if err != nil {
		// 模型未加载返回 503，其他错误返回 500
		code := 500
		msg := fmt.Sprintf("prediction failed: %v", err)
		if err.Error() == "预测模型未加载" {
			code = 503
			msg = "prediction service unavailable: model not loaded"
		}
		c.JSON(code, map[string]interface{}{
			"code":    code,
			"message": msg,
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// PredictStream 处理 SSE /api/v1/predict/:tank_id/stream
// 每30秒推送一次新的预测结果
func (h *Handler) PredictStream(_ context.Context, c *app.RequestContext) {
	tankID := c.Param("tank_id")
	if tankID == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tank_id is required",
		})
		return
	}

	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("Transfer-Encoding", "chunked")

	horizon := 12
	if v, err := strconv.Atoi(c.Query("horizon")); err == nil && v > 0 && v <= 60 {
		horizon = v
	}

	// 初始预测
	resp, err := h.predictSvc.Predict(context.Background(), tankID, horizon)
	if err != nil {
		c.WriteString(fmt.Sprintf("event: error\ndata: %s\n\n", fmt.Sprintf("prediction failed: %v", err)))
		c.Flush()
		return
	}

	// 发送初始结果
	c.WriteString(fmt.Sprintf("data: %s\n\n", mustJSON(resp)))
	c.Flush()

	// 保持连接，定期推送（简化实现：只发初始结果）
	// 实际生产中可以用 ticker 做持续推送
}

// PredictStatus 处理 GET /api/v1/predict/status
func (h *Handler) PredictStatus(_ context.Context, c *app.RequestContext) {
	loaded, version, trainedAt := h.predictSvc.ModelStatus()
	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"model_loaded": loaded,
			"version":      version,
			"trained_at":    trainedAt,
		},
	})
}

func (h *Handler) analyseStream(c *app.RequestContext, tankID string, sensorData []*domain.SensorData) {
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("Transfer-Encoding", "chunked")

	stream, err := h.analyseSvc.AnalyseStream(context.Background(), tankID, sensorData)
	if err != nil {
		c.WriteString(fmt.Sprintf("event: error\ndata: %s\n\n", fmt.Sprintf("analyse failed: %v", err)))
		c.Flush()
		return
	}

	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(msg)
		c.WriteString(fmt.Sprintf("data: %s\n\n", msg.Content))
		c.Flush()
	}
}

func convertToSensorData(dtoList []domain.SensorDataDTO) []*domain.SensorData {
	result := make([]*domain.SensorData, len(dtoList))
	for i, dto := range dtoList {
		result[i] = &domain.SensorData{
			DeviceID:    dto.DeviceID,
			TankID:      dto.TankID,
			Temperature: dto.Temperature,
			PH:          dto.PH,
			Oxygen:      dto.Oxygen,
			Ammonia:     dto.Ammonia,
			WaterLevel:  dto.WaterLevel,
			TDS:         dto.TDS,
			Nitrate:     dto.Nitrate,
			Nitrite:     dto.Nitrite,
			Chloride:    dto.Chloride,
			Timestamp:   dto.Timestamp,
		}
	}
	return result
}

// mustJSON 将对象序列化为 JSON 字符串，失败时返回错误文本
func mustJSON(v interface{}) string {
	// 使用简洁的 JSON 序列化
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error":"%v"}`, err)
	}
	return string(data)
}
