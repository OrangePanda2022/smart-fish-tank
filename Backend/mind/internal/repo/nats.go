package repo

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mind/internal/domain"

	"github.com/nats-io/nats.go"
)

// ---- Tank NATS 适配器 ----

type TankNATSRepo struct {
	nc *nats.Conn
}

func NewTankNATSRepo(url string) (*TankNATSRepo, error) {
	nc, err := nats.Connect(url,
		nats.Name("mind-tank-repo"),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("TankNATSRepo 断开连接: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("TankNATSRepo 重新连接: %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("NATS 连接失败 (tank): %w", err)
	}
	log.Printf("TankNATSRepo 已连接到 %s", url)
	return &TankNATSRepo{nc: nc}, nil
}

func (r *TankNATSRepo) GetTank(tankID string) (*domain.Tank, error) {
	req := map[string]any{
		"tank_id": tankID,
		"limit":   1, // tank 服务要求 limit 必须为 1
	}
	reqBytes, _ := json.Marshal(req)

	msg, err := r.nc.Request("tank.get", reqBytes, 3*time.Second)
	if err != nil {
		if err == nats.ErrTimeout {
			return nil, fmt.Errorf("tank 服务 NATS 超时")
		}
		return nil, fmt.Errorf("tank NATS 请求失败: %w", err)
	}

	// 先检查是否是错误响应
	var errResp map[string]string
	if json.Unmarshal(msg.Data, &errResp) == nil {
		if errMsg, ok := errResp["error"]; ok {
			return nil, fmt.Errorf("tank 服务返回错误: %s", errMsg)
		}
	}

	// 解析 tank 服务的响应
	// 注意：tank 服务的 domain.Tank 无 json tag，Go 默认按字段名（PascalCase）序列化，
	// 故这里用 PascalCase 匹配；gorm.Model 的 ID/UserID/UpdatedAt/DeletedAt 多余字段被默认忽略。
	var raw struct {
		TankID     string `json:"TankID"`
		TankName   string `json:"TankName"`
		TankSize   int    `json:"TankSize"`
		FishCount  int    `json:"FishCount"`
		FishStatus string `json:"FishStatus"`
		CreatedAt  string `json:"CreatedAt"`
	}
	if err := json.Unmarshal(msg.Data, &raw); err != nil {
		return nil, fmt.Errorf("解析 tank 响应失败: %w", err)
	}

	return &domain.Tank{
		TankID:     raw.TankID,
		TankName:   raw.TankName,
		TankSize:   raw.TankSize,
		FishCount:  raw.FishCount,
		FishStatus: raw.FishStatus,
		CreatedAt:  raw.CreatedAt,
	}, nil
}

func (r *TankNATSRepo) Close() {
	if r.nc != nil {
		r.nc.Drain()
	}
}

// ---- Sensor NATS 适配器 ----

type SensorNATSRepo struct {
	nc *nats.Conn
}

func NewSensorNATSRepo(url string) (*SensorNATSRepo, error) {
	nc, err := nats.Connect(url,
		nats.Name("mind-sensor-repo"),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("SensorNATSRepo 断开连接: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("SensorNATSRepo 重新连接: %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("NATS 连接失败 (sensor): %w", err)
	}
	log.Printf("SensorNATSRepo 已连接到 %s", url)
	return &SensorNATSRepo{nc: nc}, nil
}

func (r *SensorNATSRepo) GetLatestSensorData(tankID string, limit int) ([]*domain.SensorData, error) {
	req := map[string]any{
		"tank_id": tankID,
		"limit":   limit,
	}
	reqBytes, _ := json.Marshal(req)

	msg, err := r.nc.Request("sensor.get", reqBytes, 3*time.Second)
	if err != nil {
		if err == nats.ErrTimeout {
			return nil, fmt.Errorf("sensor 服务 NATS 超时")
		}
		return nil, fmt.Errorf("sensor NATS 请求失败: %w", err)
	}

	// 先检查是否是错误响应
	var errResp map[string]string
	if json.Unmarshal(msg.Data, &errResp) == nil {
		if errMsg, ok := errResp["error"]; ok {
			return nil, fmt.Errorf("sensor 服务返回错误: %s", errMsg)
		}
	}

	// 解析 sensor 服务的响应（[]SensorData，snake_case JSON tags）
	var rawList []struct {
		DeviceID    string  `json:"device_id"`
		TankID      string  `json:"tank_id"`
		Temperature float64 `json:"temperature"`
		PH          float64 `json:"ph"`
		Oxygen      float64 `json:"oxygen"`
		Ammonia     float64 `json:"ammonia"`
		WaterLevel  float64 `json:"water_level"`
		TDS         float64 `json:"tds"`
		Nitrate     float64 `json:"nitrate"`
		Nitrite     float64 `json:"nitrite"`
		Chloride    float64 `json:"chloride"`
		Timestamp   string  `json:"timestamp"`
	}
	if err := json.Unmarshal(msg.Data, &rawList); err != nil {
		return nil, fmt.Errorf("解析 sensor 响应失败: %w", err)
	}

	result := make([]*domain.SensorData, len(rawList))
	for i, raw := range rawList {
		result[i] = &domain.SensorData{
			DeviceID:    raw.DeviceID,
			TankID:      raw.TankID,
			Temperature: raw.Temperature,
			PH:          raw.PH,
			Oxygen:      raw.Oxygen,
			Ammonia:     raw.Ammonia,
			WaterLevel:  raw.WaterLevel,
			TDS:         raw.TDS,
			Nitrate:     raw.Nitrate,
			Nitrite:    raw.Nitrite,
			Chloride:   raw.Chloride,
			Timestamp:   raw.Timestamp,
		}
	}
	return result, nil
}

func (r *SensorNATSRepo) Close() {
	if r.nc != nil {
		r.nc.Drain()
	}
}

// Conn 返回底层 NATS 连接，供外部订阅使用
func (r *SensorNATSRepo) Conn() *nats.Conn {
	return r.nc
}
