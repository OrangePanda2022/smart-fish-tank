package model

// SensorDataDTO 表示API的传感器数据传输对象
type SensorDataDTO struct {
	DeviceID    string  `json:"deviceId"`    // 设备ID
	TankID      string  `json:"tankId"`      // 鱼缸ID
	Temperature float64 `json:"temperature"` // 温度
	PH          float64 `json:"ph"`          // pH值
	Oxygen      float64 `json:"oxygen"`      // 溶解氧
	Ammonia     float64 `json:"ammonia"`     // 氨氮含量
	WaterLevel  float64 `json:"waterLevel"`  // 水位
	Timestamp   string  `json:"timestamp"`   // 时间戳
}

// SensorDataListDTO 表示传感器数据列表
type SensorDataListDTO struct {
	Data  []SensorDataDTO `json:"data"`  // 数据列表
	Count int             `json:"count"` // 数据条数
}

// ToDTO 将领域 SensorData 转换为 DTO
func ToDTO(deviceID, tankID string, temperature, ph, oxygen, ammonia, waterLevel float64, timestamp string) SensorDataDTO {
	return SensorDataDTO{
		DeviceID:    deviceID,
		TankID:      tankID,
		Temperature: temperature,
		PH:          ph,
		Oxygen:      oxygen,
		Ammonia:     ammonia,
		WaterLevel:  waterLevel,
		Timestamp:   timestamp,
	}
}
