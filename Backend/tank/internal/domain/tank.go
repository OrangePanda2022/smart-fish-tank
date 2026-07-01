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
	TankID     string     `json:"tank_id"     gorm:"primaryKey;not null"`
	UserID     string     `json:"user_id"     gorm:"not null;index"`
	TankName   string     `json:"tank_name"   gorm:"type:varchar(100);not null"`
	TankSize   int        `json:"tank_size"   gorm:"column:tank_size;not null"`
	FishCount  int        `json:"fish_count"  gorm:"default:0"`
	FishStatus FishStatus `json:"fish_status" gorm:"type:varchar(20);default:'good'"`
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
