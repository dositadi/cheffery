package config

import (
	"fmt"
	"os"

	"github.com/dositadi/cheffery/services/shared/logger"
)

// @Description This is the redis config struct,
// @Description It loads all the configurations redis needs to start running
type RedisConfig struct {
	Protocol       string
	Username       string
	Password       string
	Host           string
	Port           int
	MaxActiveConns int
	MaxIdleConns   int
}

// @Summary load configurations for redis
// @Description load environment configurations for redis
// @Param logger
// @Produce *RedisConfig
func loadRedisConfig(logger logger.Logger) *RedisConfig {
	cfg := &RedisConfig{
		Protocol:       getStringEnvOrDefault(string(rEDIS_PROTOCOL), "redis"),
		Username:       getStringEnvOrDefault(string(rEDIS_USERNAME), "default"),
		Host:           getStringEnvOrDefault(string(rEDIS_HOST), "machine-adequate-crayon-70404.db.redis.io"),
		Port:           getIntEnvOrDefault(string(rEDIS_PORT), 18608),
		MaxActiveConns: getIntEnvOrDefault(rEDIS_MAX_ACTIVE_CONN, 30),
		MaxIdleConns:   getIntEnvOrDefault(rEDIS_MAX_IDLE_CONN, 10),
	}

	cfg.Password = os.Getenv(string(rEDIS_PASSWORD))
	if cfg.Password == "" {
		logger.PrintFatal(nil, "redis-config-loader", fmt.Sprintf("%s environment variable is required", rEDIS_PASSWORD), map[string]string{
			"Context": "config.loadRedisConfig()",
		})
		return nil
	}

	return cfg
}
