package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"sensor/internal/domain"
	"sensor/internal/repo"
)

var (
	ErrInvalidDevice = errors.New("无效的设备")
	ErrInvalidData   = errors.New("无效的传感器数据")
	ErrNotFound      = errors.New("未找到数据")
)

// SensorService 表示传感器服务
type SensorService struct {
	sensorRepo *repo.SensorRepo // 传感器数据仓储
	// deviceRepo *repo.DeviceRepo // 设备仓储
}

// NewSensorService 创建传感器服务
func NewSensorService(
	sensorRepo *repo.SensorRepo,
	// deviceRepo *repo.DeviceRepo,
) *SensorService {
	return &SensorService{
		sensorRepo: sensorRepo,
		// deviceRepo: deviceRepo,
	}
}

// ProcessMQTTMessage 处理接收到的MQTT消息
func (s *SensorService) ProcessMQTTMessage(payload []byte) error {
	// 解析JSON数据
	var data domain.SensorData
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("解析MQTT消息失败: %v", err)
		return err
	}

	// 验证设备是否有效
	// valid, err := s.deviceRepo.IsDeviceValid(data.DeviceID)
	// if err != nil {
	// 	log.Printf("验证设备失败: %v", err)
	// 	return err
	// }
	// if !valid {
	// 	log.Printf("无效的设备: %s", data.DeviceID)
	// 	return ErrInvalidDevice
	// }

	// 验证传感器数据
	// if err := data.Validate(); err != nil {
	// 	log.Printf("无效的传感器数据: %v", err)
	// 	return ErrInvalidData
	// }

	// 如果时间戳为零，使用当前时间
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}

	// 写入数据库
	if err := s.sensorRepo.Write(&data); err != nil {
		log.Printf("写入传感器数据失败: %v", err)
		return err
	}

	log.Printf("成功处理设备 %s 的传感器数据", data.DeviceID)
	return nil
}

// GetLatestByTank 获取鱼缸的最新传感器数据
func (s *SensorService) GetLatestByTank(tankID string) (*domain.SensorData, error) {
	if tankID == "" {
		return nil, ErrInvalidData
	}

	data, err := s.sensorRepo.QueryLatestByTank(tankID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrNotFound
	}

	return data, nil
}

// GetHistoryByTank 获取鱼缸的历史传感器数据
func (s *SensorService) GetHistoryByTank(tankID, start, end string, limit int) ([]domain.SensorData, error) {
	if tankID == "" {
		return nil, ErrInvalidData
	}

	// 设置默认值
	if start == "" {
		start = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}
	if end == "" {
		end = time.Now().Format(time.RFC3339)
	}
	if limit <= 0 {
		limit = 100
	}

	return s.sensorRepo.QueryHistoryByTank(tankID, start, end, limit)
}

// GetByDevice 获取特定设备的传感器数据
func (s *SensorService) GetByDevice(deviceID string, limit int) ([]domain.SensorData, error) {
	if deviceID == "" {
		return nil, ErrInvalidData
	}

	if limit <= 0 {
		limit = 100
	}

	return s.sensorRepo.QueryByDevice(deviceID, limit)
}
