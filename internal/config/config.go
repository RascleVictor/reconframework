package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App struct {
		Name        string `koanf:"name"`
		Environment string `koanf:"environment"`
		Version     string `koanf:"version"`
	} `koanf:"app"`

	Server struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	} `koanf:"server"`

	Plugins struct {
		Directory     string `koanf:"directory"`
		AutoDiscovery bool   `koanf:"auto_discovery"`
	} `koanf:"plugins"`

	Logging struct {
		Level  string `koanf:"level"`
		Format string `koanf:"format"`
	} `koanf:"logging"`
}

func Load(path string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	var cfg Config

	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal configuration: %w", err)
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return nil, fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Plugins.Directory == "" {
		return nil, fmt.Errorf("plugin directory is required")
	}

	return &cfg, nil
}
