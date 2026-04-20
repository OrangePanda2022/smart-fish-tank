package repo

import (
	"sensor/internal/database"
	"sensor/internal/domain"
)

// SensorRepo 表示传感器数据仓储
type SensorRepo struct {
	influxDB *database.InfluxDB
}

// NewSensorRepo 创建传感器数据仓储
func NewSensorRepo(influxDB *database.InfluxDB) *SensorRepo {
	return &SensorRepo{influxDB: influxDB}
}

// Write 将传感器数据写入数据库
func (r *SensorRepo) Write(data *domain.SensorData) error {
	return r.influxDB.WriteSensorData(data)
}

// QueryLatestByTank 查询鱼缸的最新传感器数据
func (r *SensorRepo) QueryLatestByTank(tankID string) (*domain.SensorData, error) {
	return r.influxDB.QueryLatestByTank(tankID)
}

// QueryHistoryByTank 查询鱼缸的历史传感器数据
func (r *SensorRepo) QueryHistoryByTank(tankID string, start, end string, limit int) ([]domain.SensorData, error) {
	return r.influxDB.QueryHistoryByTankString(tankID, start, end, limit)
}

// QueryByDevice 查询特定设备的传感器数据
func (r *SensorRepo) QueryByDevice(deviceID string, limit int) ([]domain.SensorData, error) {
	return r.influxDB.QueryByDevice(deviceID, limit)
}
