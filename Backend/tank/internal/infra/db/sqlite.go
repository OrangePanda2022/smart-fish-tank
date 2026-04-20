package db

import (
	"fmt"
	"tank/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSQLite 连接 SQLite 并自动迁移用户表
func NewSQLite(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite failed: %w", err)
	}
	if err := db.AutoMigrate(&domain.Tank{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}
	return db, nil
}
