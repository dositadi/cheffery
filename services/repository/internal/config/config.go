package config

import (
	"os"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type ApplicationConfig struct {
	RetryCfg    *RetryConfig
	PostgresEnv *PostgresEnv
}

func LoadAppConfig(logger logger.Logger) *ApplicationConfig {
	cfg := &ApplicationConfig{
		RetryCfg:    defaultRetryConfig(),
		PostgresEnv: loadPosgresEnv(logger),
	}

	return cfg
}

func getStringOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

/* func getDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if dur, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return dur
	}
	return fallback
}

func getIntOrDefault(key string, fallback int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}
	return fallback
}
*/
