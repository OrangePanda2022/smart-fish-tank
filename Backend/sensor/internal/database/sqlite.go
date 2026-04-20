package database

import (
	"fmt"
	"log"

	"sensor/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteDB 表示SQLite数据库连接
type SQLiteDB struct {
	db *gorm.DB
}

// NewSQLiteDB 创建新的SQLite数据库连接
func NewSQLiteDB(dsn string) (*SQLiteDB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接SQLite失败: %w", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&model.Device{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("SQLite数据库连接成功并完成迁移")
	return &SQLiteDB{db: db}, nil
}

// GetDB 返回底层的GORM数据库
func (s *SQLiteDB) GetDB() *gorm.DB {
	return s.db
}

// Close 关闭数据库连接
func (s *SQLiteDB) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
