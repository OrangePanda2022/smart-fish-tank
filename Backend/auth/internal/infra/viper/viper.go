package viper

import (
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func SetupViper(endpoint, path, configType string, defaults map[string]interface{}) *viper.Viper {
	v := viper.New()

	// 设置默认值
	for k, val := range defaults {
		v.SetDefault(k, val)
	}

	// 初始化熔断器配置
	settings := gobreaker.Settings{
		Name:        "etcd-config-reader",
		MaxRequests: 3,                // 半开状态下允许通过的请求数
		Interval:    5 * time.Second,  // 统计周期
		Timeout:     10 * time.Second, // 熔断器开启后，多久进入半开状态
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 如果连续失败超过 3 次，触发熔断
			return counts.ConsecutiveFailures > 3
		},
	}
	cb := gobreaker.NewCircuitBreaker(settings)

	// 配置远程提供者
	err := v.AddRemoteProvider("etcd3", endpoint, path)
	if err != nil {
		slog.Error("Failed to add remote provider",
			slog.String("endpoint", endpoint),
			slog.String("path", path),
			slog.Any("error", err),
		)
		return v
	}
	v.SetConfigType(configType)

	// 使用熔断器执行远程读取
	_, err = cb.Execute(func() (interface{}, error) {
		return nil, v.ReadRemoteConfig()
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			slog.Warn("Circuit Breaker is OPEN, skipping etcd request",
				slog.String("path", path))
		} else {
			slog.Error("Failed to read remote config from etcd",
				slog.String("path", path),
				slog.Any("error", err),
			)
		}
		slog.Info("Fallback: using default/environment configuration")
	} else {
		slog.Info("Successfully loaded configuration from etcd", slog.String("path", path))
	}

	return v
}
