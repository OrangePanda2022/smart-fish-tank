package domain

type SensorData struct {
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

type SensorDataDTO struct {
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

type SensorDataListDTO struct {
	Data  []SensorDataDTO `json:"data"`
	Count int             `json:"count"`
}
