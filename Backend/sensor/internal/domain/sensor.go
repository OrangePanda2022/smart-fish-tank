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
	TDS         float64   `json:"tds"`         // 总溶解固体 ppm
	Nitrate     float64   `json:"nitrate"`     // 硝酸根 mg/L
	Nitrite     float64   `json:"nitrite"`     // 亚硝酸根 mg/L
	Chloride    float64   `json:"chloride"`    // 氯离子 mg/L
	Timestamp   time.Time `json:"timestamp"`   // 时间戳
}

// Validate 验证传感器数据字段
// TankID 不能为空；DeviceID 允许为空（手动录入时无需指定设备）
func (s *SensorData) Validate() error {
	if s.TankID == "" {
		return ErrInvalidTankID
	}
	if s.DeviceID == "" {
		s.DeviceID = "manual-input" // 手动录入时自动填充默认设备ID
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
	if s.TDS < 0 || s.TDS > 2000 {
		return ErrInvalidTDS
	}
	if s.Nitrate < 0 || s.Nitrate > 200 {
		return ErrInvalidNitrate
	}
	if s.Nitrite < 0 || s.Nitrite > 10 {
		return ErrInvalidNitrite
	}
	if s.Chloride < 0 || s.Chloride > 1000 {
		return ErrInvalidChloride
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
	ErrInvalidTDS         = &ValidationError{Field: "tds", Message: "TDS超出范围 (0 到 2000)"}
	ErrInvalidNitrate     = &ValidationError{Field: "nitrate", Message: "硝酸根超出范围 (0 到 200)"}
	ErrInvalidNitrite     = &ValidationError{Field: "nitrite", Message: "亚硝酸根超出范围 (0 到 10)"}
	ErrInvalidChloride    = &ValidationError{Field: "chloride", Message: "氯离子超出范围 (0 到 1000)"}
)

// ValidationError 表示验证错误
type ValidationError struct {
	Field   string // 字段名
	Message string // 错误消息
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
