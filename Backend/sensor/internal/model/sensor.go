package model

import "sensor/internal/domain"

// SensorDataDTO 表示API的传感器数据传输对象
type SensorDataDTO struct {
	DeviceID    string  `json:"deviceId"`    // 设备ID
	TankID      string  `json:"tankId"`      // 鱼缸ID
	Temperature float64 `json:"temperature"` // 温度
	PH          float64 `json:"ph"`          // pH值
	Oxygen      float64 `json:"oxygen"`      // 溶解氧
	Ammonia     float64 `json:"ammonia"`     // 氨氮含量
	WaterLevel  float64 `json:"waterLevel"`  // 水位
	TDS         float64 `json:"tds"`         // 总溶解固体 ppm
	Nitrate     float64 `json:"nitrate"`     // 硝酸根 mg/L
	Nitrite     float64 `json:"nitrite"`     // 亚硝酸根 mg/L
	Chloride    float64 `json:"chlorideIon"` // 氯离子 mg/L
	Timestamp   string  `json:"timestamp"`   // 时间戳
}

// SensorDataListDTO 表示传感器数据列表
type SensorDataListDTO struct {
	Data  []SensorDataDTO `json:"data"`  // 数据列表
	Count int             `json:"count"` // 数据条数
}

// ToDTO 将领域 SensorData 转换为 DTO
func ToDTO(data *domain.SensorData) SensorDataDTO {
	return SensorDataDTO{
		DeviceID:    data.DeviceID,
		TankID:      data.TankID,
		Temperature: data.Temperature,
		PH:          data.PH,
		Oxygen:      data.Oxygen,
		Ammonia:     data.Ammonia,
		WaterLevel:  data.WaterLevel,
		TDS:         data.TDS,
		Nitrate:     data.Nitrate,
		Nitrite:     data.Nitrite,
		Chloride:    data.Chloride,
		Timestamp:   data.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToDomain 将 DTO 转换为领域 SensorData
func ToDomain(dto *SensorDataDTO) domain.SensorData {
	return domain.SensorData{
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
	}
}
