package model

import (
	"time"
)

// Device 表示系统中的IoT设备
type Device struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                 // 主键ID
	DeviceID  string    `gorm:"uniqueIndex;not null" json:"deviceId"` // 设备唯一标识
	TankID    string    `gorm:"index;not null" json:"tankId"`         // 所属鱼缸ID
	Name      string    `json:"name"`                                 // 设备名称
	Status    string    `gorm:"default:'active'" json:"status"`       // 设备状态
	CreatedAt time.Time `json:"createdAt"`                            // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`                            // 更新时间
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}
