package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 表示应用程序配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`   // 服务器配置
	Database DatabaseConfig `yaml:"database"` // 数据库配置
	InfluxDB InfluxDBConfig `yaml:"influxdb"` // InfluxDB配置
	MQTT     MQTTConfig     `yaml:"mqtt"`
	NATSURL  string         `yaml:"nats_url"` // NATS URL
}

// ServerConfig 表示HTTP服务器配置
type ServerConfig struct {
	Host string `yaml:"host"` // 监听地址
	Port int    `yaml:"port"` // 监听端口
}

// DatabaseConfig 表示SQLite数据库配置
type DatabaseConfig struct {
	Driver string `yaml:"driver"` // 数据库驱动
	DSN    string `yaml:"dsn"`    // 数据源名称
}

// InfluxDBConfig 表示InfluxDB配置
type InfluxDBConfig struct {
	URL    string `yaml:"url"`    // InfluxDB URL
	Token  string `yaml:"token"`  // 认证令牌
	Org    string `yaml:"org"`    // 组织名称
	Bucket string `yaml:"bucket"` // 存储桶名称
}

// MQTTConfig 表示MQTT配置
type MQTTConfig struct {
	Broker   string `yaml:"broker"`    // MQTT broker地址
	ClientID string `yaml:"client_id"` // 客户端ID
	Username string `yaml:"username"`  // 用户名
	Password string `yaml:"password"`  // 密码
	Topic    string `yaml:"topic"`     // 订阅主题
	QOS      int    `yaml:"qos"`       // QoS等级
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

	return &cfg, nil
}
