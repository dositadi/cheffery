package config

import (
	"fmt"
	"os"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type RedisConfig struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     int
}

func loadRedisConfig(logger logger.Logger) *RedisConfig {
	cfg := &RedisConfig{
		Protocol: getStringEnvOrDefault(string(rEDIS_PROTOCOL), "redis"),
		Username: getStringEnvOrDefault(string(rEDIS_USERNAME), "default"),
		Host:     getStringEnvOrDefault(string(rEDIS_HOST), "machine-adequate-crayon-70404.db.redis.io"),
		Port:     getIntEnvOrDefault(string(rEDIS_PORT), 18608),
	}

	cfg.Password = os.Getenv(string(rEDIS_PASSWORD))
	if cfg.Password == "" {
		logger.PrintFatal(nil, fmt.Sprintf("%s environment variable is required", rEDIS_PASSWORD), map[string]string{
			"Context": "config.loadRedisConfig()",
		})
		return nil
	}

	return cfg
}
