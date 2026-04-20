package repo

import (
	"context"
	"errors"
	"fmt"
	"tank/internal/domain"

	"gorm.io/gorm"
)

type TankRepo struct {
	db *gorm.DB
}

func NewTankRepo(db *gorm.DB) *TankRepo {
	return &TankRepo{db: db}
}

func (t *TankRepo) Create(ctx context.Context, tank *domain.Tank) error {
	if err := t.db.WithContext(ctx).Create(tank).Error; err != nil {
		return fmt.Errorf("create tank failed: %w", err)
	}
	return nil
}

func (t *TankRepo) Delete(ctx context.Context, tankID string) error {
	res := t.db.WithContext(ctx).Where("tank_id = ?", tankID).Delete(&domain.Tank{})
	if res.Error != nil {
		return fmt.Errorf("soft delete tank failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (t *TankRepo) GetTankByTankID(ctx context.Context, tankID string) (*domain.Tank, error) {
	var tank domain.Tank
	err := t.db.WithContext(ctx).Where("tank_id = ?", tankID).First(&tank).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tank by id failed: %w", err)
	}
	return &tank, nil
}

func (t *TankRepo) GetTanksByUserID(ctx context.Context, userID string) ([]*domain.Tank, error) {
	var tanks []*domain.Tank

	err := t.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&tanks).Error

	if err != nil {
		return nil, err
	}

	return tanks, nil
}

func (t *TankRepo) UpdateTank(ctx context.Context, tankID string, tank *domain.Tank) error {
	res := t.db.WithContext(ctx).Model(&domain.Tank{}).
		Where("tank_id = ?", tankID).
		Updates(tank)
	if res.Error != nil {
		return fmt.Errorf("update tank failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
