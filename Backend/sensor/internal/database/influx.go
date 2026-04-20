package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"sensor/internal/config"
	"sensor/internal/domain"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

// InfluxDB 表示时序数据库客户端
type InfluxDB struct {
	client *influxdb3.Client
	org    string
	bucket string
	// influx配置
	mu   sync.RWMutex
	data map[string][]domain.SensorData
}

// NewInfluxDB 创建新的InfluxDB连接
func NewInfluxDB(cfg config.InfluxDBConfig) (*InfluxDB, error) {
	log.Printf("连接到InfluxDB: %s (bucket: %s, org: %s)", cfg.URL, cfg.Bucket, cfg.Org)

	// 创建InfluxDB客户端配置
	clientConfig := influxdb3.ClientConfig{
		Host:         cfg.URL,
		Token:        cfg.Token,
		Organization: cfg.Org,
		Database:     cfg.Bucket,
	}

	// 创建客户端
	client, err := influxdb3.New(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建InfluxDB客户端失败: %w", err)
	}

	return &InfluxDB{
		client: client,
		org:    cfg.Org,
		bucket: cfg.Bucket,
		data:   make(map[string][]domain.SensorData),
	}, nil
}

// WriteSensorData 将传感器数据写入InfluxDB
func (i *InfluxDB) WriteSensorData(data *domain.SensorData) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 首先写入内存缓存（用于快速查询）
	key := data.TankID
	if _, ok := i.data[key]; !ok {
		i.data[key] = []domain.SensorData{}
	}
	i.data[key] = append(i.data[key], *data)
	fmt.Println(key, data)

	// 同时按设备ID存储，方便设备查询
	deviceKey := "device:" + data.DeviceID
	if _, ok := i.data[deviceKey]; !ok {
		i.data[deviceKey] = []domain.SensorData{}
	}
	i.data[deviceKey] = append(i.data[deviceKey], *data)

	// 写入InfluxDB
	// 创建数据点
	point := influxdb3.NewPoint(
		"sensor_data",
		map[string]string{
			"device_id": data.DeviceID,
			"tank_id":   data.TankID,
		},
		map[string]interface{}{
			"temperature": data.Temperature,
			"ph":          data.PH,
			"oxygen":      data.Oxygen,
			"ammonia":     data.Ammonia,
			"water_level": data.WaterLevel,
		},
		data.Timestamp,
	)

	// 写入InfluxDB
	ctx := context.Background()
	err := i.client.WritePoints(ctx, []*influxdb3.Point{point})
	if err != nil {
		log.Printf("写入InfluxDB失败: %v (数据已保存到内存缓存)", err)
		return nil // 返回nil，因为数据已经保存到内存缓存中
	}

	log.Printf("已写入传感器数据到InfluxDB - 鱼缸: %s, 设备: %s", data.TankID, data.DeviceID)
	return nil
}

// QueryLatestByTank 查询鱼缸的最新传感器数据
func (i *InfluxDB) QueryLatestByTank(tankID string) (*domain.SensorData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// 先尝试从内存缓存获取
	records, ok := i.data[tankID]
	fmt.Println(tankID, records, ok)
	if ok {
		// 获取最新的一条记录
		latest := records[len(records)-1]
		return &latest, nil
	}

	// 内存缓存中没有，从InfluxDB查询
	query := fmt.Sprintf(`
		SELECT *
		FROM sensor_data
		WHERE tank_id = '%s'
		ORDER BY time DESC
		LIMIT 1
	`, tankID)

	iter, err := i.client.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %w", err)
	}
	defer func() {
		// QueryIterator 没有 Close 方法，直接忽略
	}()

	if !iter.Next() {
		return nil, fmt.Errorf("未找到鱼缸 %s 的数据", tankID)
	}

	// 解析查询结果
	row := iter.Value()
	timestamp := time.Time{}
	if t, ok := row["time"]; ok {
		switch v := t.(type) {
		case time.Time:
			timestamp = v
		}
	}

	return &domain.SensorData{
		DeviceID:    toString(row["device_id"]),
		TankID:      tankID,
		Temperature: toFloat64(row["temperature"]),
		PH:          toFloat64(row["ph"]),
		Oxygen:      toFloat64(row["oxygen"]),
		Ammonia:     toFloat64(row["ammonia"]),
		WaterLevel:  toFloat64(row["water_level"]),
		Timestamp:   timestamp,
	}, nil
}

// QueryHistoryByTank 查询鱼缸的历史传感器数据
func (i *InfluxDB) QueryHistoryByTank(tankID string, start, end time.Time, limit int) ([]domain.SensorData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// 先尝试从内存缓存获取
	records, ok := i.data[tankID]
	if ok && len(records) > 0 {
		// 按时间范围过滤
		var result []domain.SensorData
		for _, r := range records {
			if (start.IsZero() || r.Timestamp.After(start) || r.Timestamp.Equal(start)) &&
				(end.IsZero() || r.Timestamp.Before(end) || r.Timestamp.Equal(end)) {
				result = append(result, r)
				if limit > 0 && len(result) >= limit {
					break
				}
			}
		}
		return result, nil
	}

	// 内存缓存中没有，从InfluxDB查询
	timeRange := ""
	if !start.IsZero() {
		timeRange += fmt.Sprintf(" AND time >= timestamp '%s'", start.Format(time.RFC3339))
	}
	if !end.IsZero() {
		timeRange += fmt.Sprintf(" AND time <= timestamp '%s'", end.Format(time.RFC3339))
	}

	limitStr := ""
	if limit > 0 {
		limitStr = fmt.Sprintf(" LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sensor_data
		WHERE tank_id = '%s'%s
		ORDER BY time DESC
		%s
	`, tankID, timeRange, limitStr)

	iter, err := i.client.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %w", err)
	}

	var result []domain.SensorData
	for iter.Next() {
		row := iter.Value()
		timestamp := time.Time{}
		if t, ok := row["time"]; ok {
			switch v := t.(type) {
			case time.Time:
				timestamp = v
			}
		}

		result = append(result, domain.SensorData{
			DeviceID:    toString(row["device_id"]),
			TankID:      tankID,
			Temperature: toFloat64(row["temperature"]),
			PH:          toFloat64(row["ph"]),
			Oxygen:      toFloat64(row["oxygen"]),
			Ammonia:     toFloat64(row["ammonia"]),
			WaterLevel:  toFloat64(row["water_level"]),
			Timestamp:   timestamp,
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("查询InfluxDB出错: %w", err)
	}

	return result, nil
}

// QueryHistoryByTankString 使用字符串时间参数查询历史传感器数据
func (i *InfluxDB) QueryHistoryByTankString(tankID, startStr, endStr string, limit int) ([]domain.SensorData, error) {
	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, fmt.Errorf("无效的开始时间: %w", err)
		}
	}
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, fmt.Errorf("无效的结束时间: %w", err)
		}
	}
	return i.QueryHistoryByTank(tankID, start, end, limit)
}

// QueryByDevice 查询特定设备的传感器数据
func (i *InfluxDB) QueryByDevice(deviceID string, limit int) ([]domain.SensorData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// 先尝试从内存缓存获取
	deviceKey := "device:" + deviceID
	records, ok := i.data[deviceKey]
	if ok && len(records) > 0 {
		// 获取最后N条记录
		start := len(records) - limit
		if start < 0 {
			start = 0
		}
		if limit > 0 && len(records) > limit {
			return records[start:], nil
		}
		return records, nil
	}

	// 内存缓存中没有，从InfluxDB查询
	limitStr := ""
	if limit > 0 {
		limitStr = fmt.Sprintf(" LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT *
		FROM sensor_data
		WHERE device_id = '%s'
		ORDER BY time DESC%s
	`, deviceID, limitStr)

	iter, err := i.client.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %w", err)
	}

	var result []domain.SensorData
	for iter.Next() {
		row := iter.Value()
		timestamp := time.Time{}
		if t, ok := row["time"]; ok {
			switch v := t.(type) {
			case time.Time:
				timestamp = v
			}
		}

		result = append(result, domain.SensorData{
			DeviceID:    deviceID,
			TankID:      toString(row["tank"]),
			Temperature: toFloat64(row["temperature"]),
			PH:          toFloat64(row["ph"]),
			Oxygen:      toFloat64(row["oxygen"]),
			Ammonia:     toFloat64(row["ammonia"]),
			WaterLevel:  toFloat64(row["water_level"]),
			Timestamp:   timestamp,
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("查询InfluxDB出错: %w", err)
	}

	return result, nil
}

// Close 关闭InfluxDB连接
func (i *InfluxDB) Close() {
	if i.client != nil {
		i.client.Close()
	}
	log.Println("InfluxDB连接已关闭")
}

// toFloat64 安全转换为float64
func toFloat64(v any) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}

// toString 安全转换为string
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
