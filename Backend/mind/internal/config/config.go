package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Milvus   MilvusConfig
	LLM      LLMConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type MilvusConfig struct {
	Addr       string
	Collection string
	UserName   string
	Password   string
	DBName     string
}

type LLMConfig struct {
	APIKey string
	Model  string
}

type DatabaseConfig struct {
	DSN string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "6788"),
		},
		Milvus: MilvusConfig{
			Addr:       getEnv("MILVUS_ADDR", "localhost:19530"),
			Collection: getEnv("MILVUS_COLLECTION", "aquaculture_knowledge"),
			UserName:   getEnv("MILVUS_USERNAME", ""),
			Password:   getEnv("MILVUS_PASSWORD", ""),
			DBName:     getEnv("MILVUS_DB_NAME", "default"),
		},
		LLM: LLMConfig{
			APIKey: getEnv("LLM_API_KEY", ""),
			Model:  getEnv("LLM_MODEL", "gpt-4"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DATABASE_DSN", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
