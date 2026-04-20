package domain

import (
	"tank/internal/util"

	"gorm.io/gorm"
)

const (
	FishGreen  FishStatus = "good"
	FishYellow FishStatus = "not good"
	FishRed    FishStatus = "bad"
)

type FishStatus string

type Tank struct {
	TankID     string     `gorm:"primaryKey;not null"`
	UserID     string     `gorm:"not null;index"`
	TankName   string     `gorm:"type:varchar(100);not null"`
	TankSize   int        `gorm:"column:tank_size;not null"`
	FishCount  int        `gorm:"default:0"`
	FishStatus FishStatus `gorm:"type:varchar(20);default:'good'"`
	gorm.Model
}

func (u *Tank) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果 ID 已经有值，则不覆盖
	if u.TankID == "" {
		uid, err := util.NewUUIDv7()
		if err != nil {
			return err
		}
		u.TankID = uid
	}
	return nil
}
