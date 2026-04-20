package service

import (
	"context"
	"fmt"
	"tank/internal/domain"
	"tank/internal/repo"
)

type TankService struct {
	tankRepo *repo.TankRepo
}

func NewTankService(tankRepo *repo.TankRepo) *TankService {
	return &TankService{tankRepo: tankRepo}
}

func (t *TankService) GetTankByTankID(ctx context.Context, userID, tankID string) (*domain.Tank, error) {
	tank, _ := t.tankRepo.GetTankByTankID(ctx, tankID)
	if tank.UserID != userID {
		return nil, fmt.Errorf("not your tank")
	}
	return tank, nil
}

func (t *TankService) GetTanksByUserID(ctx context.Context, userID string) ([]*domain.Tank, error) {
	return t.tankRepo.GetTanksByUserID(ctx, userID)
}

func (t *TankService) CreateTank(ctx context.Context, userID string, tankName string, tankSize int) error {
	tank := &domain.Tank{
		UserID:   userID,
		TankName: tankName,
		TankSize: tankSize,
	}
	return t.tankRepo.Create(ctx, tank)
}

func (t *TankService) DeleteTank(ctx context.Context, userID, tankID string) error {
	tank, _ := t.tankRepo.GetTankByTankID(ctx, tankID)
	if tank.UserID != userID {
		return fmt.Errorf("not your tank")
	}
	return t.tankRepo.Delete(ctx, tankID)
}

func (t *TankService) UpdateTankByTankID(ctx context.Context, userID, tankID string, tankName string, tankSize int) error {
	if _, err := t.GetTankByTankID(ctx, userID, tankID); err != nil {
		return err
	}
	tank := &domain.Tank{
		TankName: tankName,
		TankSize: tankSize,
	}
	return t.tankRepo.UpdateTank(ctx, tankID, tank)
}
