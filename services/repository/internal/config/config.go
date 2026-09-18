package config

import (
	"os"
	"strconv"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type ApplicationConfig struct {
	RetryCfg    *RetryConfig
	PostgresEnv *PostgresEnv
	ServerPort  uint16
}

func LoadAppConfig(logger logger.Logger) *ApplicationConfig {
	cfg := &ApplicationConfig{
		RetryCfg:    defaultRetryConfig(),
		PostgresEnv: loadPosgresEnv(logger),
		ServerPort:  getIntOrDefault("REPO_GRPC_PORT", 50052),
	}

	return cfg
}

func getStringOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if dur, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return dur
	}
	return fallback
}

func getIntOrDefault(key string, fallback uint16) uint16 {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		if val > 0 {
			return uint16(val)
		}
	}
	return fallback
}
