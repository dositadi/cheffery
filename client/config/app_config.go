package config

import (
	"os"
	"strconv"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type AppConfig struct {
	Service *ServiceConfig
	Server  *ServerConfig
}

func LoadAppConfig(logger logger.Logger) *AppConfig {
	cfg := &AppConfig{
		Service: loadServiceConfig(logger),
		Server:  loadServerConfig(),
	}
	return cfg
}

func getStringOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getIntOrDefault(key string, fallback uint16) uint16 {
	if val, err := strconv.Atoi(key); err == nil && val > 0 {
		return uint16(val)
	}
	return fallback
}

func getDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if dur, err := time.ParseDuration(key); err == nil {
		return dur
	}
	return fallback
}
