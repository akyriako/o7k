package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultCacheTTL        = 300
	defaultRefreshInterval = 30
)

type Config struct {
	Cache   Cache   `yaml:"cache"`
	Log     Log     `yaml:"log"`
	Refresh Refresh `yaml:"refresh"`
}

type Cache struct {
	Enabled         bool     `yaml:"enabled"`
	TTL             int      `yaml:"ttl"`
	CleanupInterval int      `yaml:"cleanup_interval"`
	Resources       []string `yaml:"resources"`
}

type Log struct {
	Level slog.Level `yaml:"level"`
}

type Refresh struct {
	Interval      int  `yaml:"interval"`
	ShowTimestamp bool `yaml:"show_timestamp"`
}

func Default() Config {
	return Config{
		Cache: Cache{
			Enabled:         true,
			TTL:             defaultCacheTTL,
			CleanupInterval: defaultCacheTTL,
			Resources: []string{
				"catalog",
				"regions",
				"services",
				"endpoints",
				"flavors",
				"images",
			},
		},
		Log: Log{
			Level: slog.LevelInfo,
		},
		Refresh: Refresh{
			Interval: defaultRefreshInterval,
		},
	}
}

func Load() (Config, error) {
	cfg := Default()

	dir, err := os.UserConfigDir()
	if err != nil {
		return cfg, err
	}

	path := filepath.Join(dir, "o7k", "config.yaml")

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if cfg.Cache.Enabled {
		if cfg.Cache.TTL <= 0 {
			cfg.Cache.TTL = defaultCacheTTL
		}

		if cfg.Cache.CleanupInterval <= 0 {
			cfg.Cache.CleanupInterval = cfg.Cache.TTL
		}
	}

	if cfg.Refresh.Interval <= 0 {
		cfg.Refresh.Interval = defaultRefreshInterval
	}

	return cfg, nil
}
