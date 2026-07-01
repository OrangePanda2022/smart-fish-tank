package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server     ServerConfig
	NATS       NATSConfig
	Milvus     MilvusConfig
	LLM        LLMConfig
	Database   DatabaseConfig
	Prediction PredictionConfig
}

type NATSConfig struct {
	URL string
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

type PredictionConfig struct {
	ModelPath   string
	Horizon     int
	StateWeight float64
	CtrlWeight  float64
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "6788"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
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
		Prediction: PredictionConfig{
			ModelPath:   getEnv("KOOPMAN_MODEL_PATH", "./models/koopman_model.json"),
			Horizon:     parseIntEnv("MPC_HORIZON", "12"),
			StateWeight: parseFloatEnv("MPC_STATE_WEIGHT", "10.0"),
			CtrlWeight:  parseFloatEnv("MPC_CTRL_WEIGHT", "1.0"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseIntEnv(key, defaultValue string) int {
	v := getEnv(key, defaultValue)
	n, err := strconv.Atoi(v)
	if err != nil {
		n, _ = strconv.Atoi(defaultValue)
	}
	return n
}

func parseFloatEnv(key, defaultValue string) float64 {
	v := getEnv(key, defaultValue)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		f, _ = strconv.ParseFloat(defaultValue, 64)
	}
	return f
}
