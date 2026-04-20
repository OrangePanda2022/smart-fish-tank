package config

import (
	"gateway/internal/domain/upstream"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type FallbackNode struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type FallbackConfig struct {
	Routes    map[string]string          `yaml:"routes"`
	Upstreams map[string][]upstream.Node `yaml:"upstreams"`
}

func LoadFallback(filePath string, logger *zap.Logger) (*FallbackConfig, error) {
	var k = koanf.New(".")

	// 加载文件
	if err := k.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		return nil, err
	}

	// 将配置映射到结构体
	var config FallbackConfig
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type RouteMapConfig struct {
	Routes map[string]string `yaml:"routes"`
}

func LoadRouteMap(filePath string) (map[string]string, error) {
	var k = koanf.New(".")

	if err := k.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		return nil, err
	}

	var config RouteMapConfig
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}
	return config.Routes, nil
}
