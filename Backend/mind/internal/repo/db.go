package repo

import (
	"mind/internal/domain"
)

// SensorRepository 抽象传感器数据访问；nats (SensorNATSRepo) 实现满足。
type SensorRepository interface {
	GetLatestSensorData(tankID string, limit int) ([]*domain.SensorData, error)
}

// TankRepository 抽象鱼缸数据访问；nats (TankNATSRepo) 实现满足。
type TankRepository interface {
	GetTank(tankID string) (*domain.Tank, error)
}
