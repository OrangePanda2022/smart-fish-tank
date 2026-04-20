package domain

import "time"

// SensorData 表示传感器数据的领域实体
type SensorData struct {
	DeviceID    string    `json:"device_id"`   // 设备ID
	TankID      string    `json:"tank_id"`     // 鱼缸ID
	Temperature float64   `json:"temperature"` // 温度
	PH          float64   `json:"ph"`          // pH值
	Oxygen      float64   `json:"oxygen"`      // 溶解氧
	Ammonia     float64   `json:"ammonia"`     // 氨氮含量
	WaterLevel  float64   `json:"water_level"` // 水位
	Timestamp   time.Time `json:"timestamp"`   // 时间戳
}

// Validate 验证传感器数据字段
func (s *SensorData) Validate() error {
	if s.DeviceID == "" {
		return ErrInvalidDeviceID
	}
	if s.TankID == "" {
		return ErrInvalidTankID
	}
	if s.Temperature < -50 || s.Temperature > 100 {
		return ErrInvalidTemperature
	}
	if s.PH < 0 || s.PH > 14 {
		return ErrInvalidPH
	}
	if s.Oxygen < 0 || s.Oxygen > 20 {
		return ErrInvalidOxygen
	}
	if s.Ammonia < 0 || s.Ammonia > 10 {
		return ErrInvalidAmmonia
	}
	if s.WaterLevel < 0 || s.WaterLevel > 100 {
		return ErrInvalidWaterLevel
	}
	return nil
}

// 错误定义
var (
	ErrInvalidDeviceID    = &ValidationError{Field: "deviceId", Message: "设备ID不能为空"}
	ErrInvalidTankID      = &ValidationError{Field: "tankId", Message: "鱼缸ID不能为空"}
	ErrInvalidTemperature = &ValidationError{Field: "temperature", Message: "温度超出范围 (-50 到 100)"}
	ErrInvalidPH          = &ValidationError{Field: "ph", Message: "pH值超出范围 (0 到 14)"}
	ErrInvalidOxygen      = &ValidationError{Field: "oxygen", Message: "溶解氧超出范围 (0 到 20)"}
	ErrInvalidAmmonia     = &ValidationError{Field: "ammonia", Message: "氨氮含量超出范围 (0 到 10)"}
	ErrInvalidWaterLevel  = &ValidationError{Field: "waterLevel", Message: "水位超出范围 (0 到 100)"}
)

// ValidationError 表示验证错误
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误消息
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
