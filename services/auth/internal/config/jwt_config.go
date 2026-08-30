package config

import (
	"os"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type JWTConfig struct {
	AccessKey  string
	RefreshKey string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func loadJWTConfig(logger logger.Logger) *JWTConfig {
	cfg := &JWTConfig{
		AccessTTL:  getDurationOrDefault("ACCESS_TTL", time.Duration(15*time.Minute)),
		RefreshTTL: getDurationOrDefault("REFRESH_TTL", time.Duration(168*time.Minute)),
	}

	cfg.AccessKey = os.Getenv("ACCESS_TOKEN_KEY")
	if cfg.AccessKey == "" {
		logger.PrintFatal(nil, "jwt-config-loader", "ACCESS_TOKEN_KEY environment variable is required", map[string]string{
			"Context": "config.loadRedisConfig()",
		})
	}

	cfg.RefreshKey = os.Getenv("REFRESH_TOKEN_KEY")
	if cfg.RefreshKey == "" {
		logger.PrintFatal(nil, "jwt-config-loader", "ACCESS_TOKEN_KEY environment variable is required", map[string]string{
			"Context": "config.loadRedisConfig()",
		})
	}

	return cfg
}
