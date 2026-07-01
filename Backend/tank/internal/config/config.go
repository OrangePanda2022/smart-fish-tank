package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示应用程序配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`   // 服务器配置
	Database DatabaseConfig `yaml:"database"` // 数据库配置
	NATSURL  string         `yaml:"nats_url"` // NATS URL
	Stream   StreamConfig   `yaml:"stream"`   // 直播流配置
}

// ServerConfig 表示HTTP服务器配置
type ServerConfig struct {
	Host string `yaml:"host"` // 监听地址
	Port int    `yaml:"port"` // 监听端口
}

// DatabaseConfig 表示SQLite数据库配置
type DatabaseConfig struct {
	DSN string `yaml:"dsn"` // 数据源名称
}

// StreamConfig 表示直播流配置
type StreamConfig struct {
	MaxRingSize   int `yaml:"max_ring_size"`   // 每个鱼缸的帧环形缓冲区容量，默认3
	ViewerBufSize int `yaml:"viewer_buf_size"` // 观看者通知channel缓冲区大小，默认1
}

// Load 从YAML文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8084
	}
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = "./tank.db"
	}
	if cfg.NATSURL == "" {
		cfg.NATSURL = "nats://localhost:4222"
	}
	if cfg.Stream.MaxRingSize == 0 {
		cfg.Stream.MaxRingSize = 3
	}
	if cfg.Stream.ViewerBufSize == 0 {
		cfg.Stream.ViewerBufSize = 1
	}

	return &cfg, nil
}
