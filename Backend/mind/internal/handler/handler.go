package handler

import (
	"context"
	"fmt"
	"mind/internal/domain"
	"mind/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Handler struct {
	analyseSvc *service.AnalyseService
}

func NewHandler(analyseSvc *service.AnalyseService) *Handler {
	return &Handler{
		analyseSvc: analyseSvc,
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
			Timestamp:   dto.Timestamp,
		}
	}
	return result
}

// func NewHandlerWithDeps(
// 	tankRepo *repo.TankRepo,
// 	sensorRepo *repo.SensorRepo,
// 	milvusAddr, milvusCollection string,
// ) *Handler {
// 	milvus := vector.NewMilvusClient(milvusAddr, milvusCollection)
// 	return NewHandler(tankRepo, sensorRepo, milvus)
// }

// var _ context.Context = (*handlerContext)(nil)

// type handlerContext struct {
// 	ctx context.Context
// }

// func (hc *handlerContext) Deadline() (time.Time, bool) {
// 	return hc.ctx.Deadline()
// }

// func (hc *handlerContext) Done() <-chan struct{} {
// 	return hc.ctx.Done()
// }

// func (hc *handlerContext) Err() error {
// 	return hc.ctx.Err()
// }

// func (hc *handlerContext) Value(key any) any {
// 	return hc.ctx.Value(key)
// }

// func wrapContext(ctx context.Context) context.Context {
// 	return &handlerContext{ctx: ctx}
// }

// func parseIntWithDefault(s string, defaultVal int) int {
// 	if v, err := strconv.Atoi(s); err == nil {
// 		return v
// 	}
// 	return defaultVal
// }
