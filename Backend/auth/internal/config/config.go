package config

import (
	"auth/internal/infra/viper"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	once   sync.Once
	global Config
)

// Config 定义应用运行配置，全部由环境变量驱动并提供安全默认值
type Config struct {
	AppName                     string
	HTTPAddr                    string
	ReadTimeout                 time.Duration
	WriteTimeout                time.Duration
	ShutdownTimeout             time.Duration
	SQLitePath                  string
	JWTIssuer                   string
	JWTAccessSecret             string
	JWTRefreshSecret            string
	AccessTTL                   time.Duration
	RefreshTTL                  time.Duration
	ResetTokenTTL               time.Duration
	EmailMapTTL                 time.Duration
	RevokeMapTTL                time.Duration
	UserByIDTTL                 time.Duration
	UserByIDNullTTL             time.Duration
	BufferSize                  int
	ClientSecret                string
	OAuthClients                map[string]string
	OAuthRedirectURIs           map[string][]string
	HydraPublicURL              string
	HydraAdminURL               string
	HydraIntrospectClientID     string
	HydraIntrospectClientSecret string
	EmailBloomKey               string
	EmailBloomErr               float64
	EmailBloomCap               uint64
	RedisAddrs                  []string
	RedisPassword               string
	RedisDB                     int
	KafkaBrokers                []string
	KafkaTopic                  string
	CasbinPath                  string
	KeyFuncURLs                 []string
	SessionCookieMaxAge         int
	HydraLoginRememberFor       int64
	HydraConsentRememberFor     int64
	HttpRateLimitRPS            int
	HttpRateLimitBurstlogger    int
}

// Load 从环境变量加载配置，保证示例开箱可运行
func Load() Config {
	once.Do(func() {
		// 1. 定义默认值映射
		defaults := map[string]interface{}{
			"AppName":          "auth-service",
			"HTTPAddr":         ":8081",
			"ReadTimeout":      10 * time.Second,
			"WriteTimeout":     15 * time.Second,
			"ShutdownTimeout":  10 * time.Second,
			"SQLitePath":       "./auth.db",
			"JWTIssuer":        "auth-service",
			"JWTAccessSecret":  "access-secret-change-me",
			"JWTRefreshSecret": "refresh-secret-change-me",
			"AccessTTL":        15 * time.Minute,
			"RefreshTTL":       7 * 24 * time.Hour,
			"ResetTokenTTL":    15 * time.Minute,
			"EmailMapTTL":      60 * time.Minute,
			"RevokeMapTTL":     15 * time.Minute,
			"UserByIDTTL":      15 * time.Minute,
			"UserByIDNullTTL":  2 * time.Minute,
			"BufferSize":       1024,
			"ClientSecret":     "test",
			"OAuthClients":     map[string]string{"5eb8cc07-aed7-4570-aa93-c98f678b3d95": "test123456"},
			"OAuthRedirectURIs": map[string][]string{
				"test": {"http://localhost:3000/callback"},
			},
			"HydraPublicURL":              "http://127.0.0.1:4444",
			"HydraAdminURL":               "http://127.0.0.1:4445",
			"HydraIntrospectClientID":     "test",
			"HydraIntrospectClientSecret": "test",
			"EmailBloomKey":               "bf:email:exists",
			"EmailBloomErr":               0.01,
			"EmailBloomCap":               100000,
			"RedisAddrs":                  []string{""},
			"RedisPassword":               "",
			"RedisDB":                     0,
			"KafkaBrokers":                []string{""},
			"KafkaTopic":                  "auth.audit",
			"CasbinPath":                  "./internal/config/rbac_model.conf",
			"KeyFuncURLs":                 []string{"http://localhost:4444/.well-known/jwks.json"},
			"SessionCookieMaxAge":         3600,
			"HydraLoginRememberFor":       int64(3600),
			"HydraConsentRememberFor":     int64(30 * 24 * 3600),
			"HttpRateLimitRPS":            20,
			"HttpRateLimitBurstlogger":    40,
		}

		// 调用 Viper Setup
		ETCDAddr := getEnv("ETCD_ADDR", "")
		ETCDPath := getEnv("ETCD_PATH", "/config/auth-service.json")
		v := viper.SetupViper(
			ETCDAddr,
			ETCDPath,
			"json",
			defaults,
		)

		// 将 Viper 里的数据解析到结构体
		if err := v.Unmarshal(&global); err != nil {
			slog.Error("Failed to unmarshal config into struct", slog.Any("error", err))
		}
	})

	return global
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getIntEnv 读取整型环境变量，非法值回退到默认值
func getIntEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func getUint64Env(key string, fallback uint64) uint64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}

func getFloatEnv(key string, fallback float64) float64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return v
}

// getDurationEnv 读取 time.ParseDuration 格式的时长环境变量
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return v
}

// splitCSV 解析逗号分隔配置，并忽略多余空白字符
func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	result := make([]string, 0, 4)
	cur := ""
	for _, ch := range raw {
		if ch == ',' {
			if cur != "" {
				result = append(result, cur)
			}
			cur = ""
			continue
		}
		if ch == ' ' || ch == '\t' || ch == '\n' {
			continue
		}
		cur += string(ch)
	}
	if cur != "" {
		result = append(result, cur)
	}
	return result
}
