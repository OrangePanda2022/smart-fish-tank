package repo

import (
	"mind/internal/domain"
	"sync"
	"time"
)

type TankRepo struct {
	store map[string]*domain.Tank
	mu    sync.RWMutex
}

func NewTankRepo() *TankRepo {
	repo := &TankRepo{store: make(map[string]*domain.Tank)}
	repo.initSampleData()
	return repo
}

func (r *TankRepo) initSampleData() {
	r.store["tank_001"] = &domain.Tank{
		TankID:      "tank-001",
		FishSpecies: []string{"孔雀鱼", "灯鱼", "黑玛丽"},
		CreatedAt:   time.Now().AddDate(0, -3, 0).Format(time.RFC3339),
	}
	r.store["tank_002"] = &domain.Tank{
		TankID:      "tank-002",
		FishSpecies: []string{"金鱼", "锦鲤"},
		CreatedAt:   time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
	}
}

func (r *TankRepo) GetTank(tankID string) (*domain.Tank, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if tank, ok := r.store[tankID]; ok {
		return tank, nil
	}
	return nil, nil
}

type SensorRepo struct {
	store map[string][]*domain.SensorData
	mu    sync.RWMutex
}

func NewSensorRepo() *SensorRepo {
	repo := &SensorRepo{store: make(map[string][]*domain.SensorData)}
	repo.initSampleData()
	return repo
}

func (r *SensorRepo) initSampleData() {
	now := time.Now()
	r.store["tank-001"] = []*domain.SensorData{
		{
			DeviceID:    "sensor-001",
			TankID:      "tank-001",
			Temperature: 28.5,
			PH:          7.2,
			Oxygen:      6.5,
			Ammonia:     0.8,
			WaterLevel:  75.0,
			Timestamp:   now.Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			DeviceID:    "sensor-001",
			TankID:      "tank-001",
			Temperature: 28.2,
			PH:          7.1,
			Oxygen:      6.8,
			Ammonia:     0.6,
			WaterLevel:  76.0,
			Timestamp:   now.Format(time.RFC3339),
		},
	}
	r.store["tank-002"] = []*domain.SensorData{
		{
			DeviceID:    "sensor-002",
			TankID:      "tank-002",
			Temperature: 24.0,
			PH:          7.4,
			Oxygen:      7.2,
			Ammonia:     0.1,
			WaterLevel:  80.0,
			Timestamp:   now.Format(time.RFC3339),
		},
	}
}

func (r *SensorRepo) GetLatestSensorData(tankID string, limit int) ([]*domain.SensorData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if data, ok := r.store[tankID]; ok {
		if limit > 0 && len(data) > limit {
			return data[len(data)-limit:], nil
		}
		return data, nil
	}
	return nil, nil
}

func (r *SensorRepo) SaveSensorData(tankID string, data []*domain.SensorData) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[tankID] = append(r.store[tankID], data...)
	return nil
}
