package config

import (
	"os"
	"strconv"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
)

// @Description This is the auth service config file,
// @Description It loads all the configurations the auth service needs to start running
type Config struct {
	RedisConfig *RedisConfig // Redis config
}

func New(logger logger.Logger) *Config {
	return loadConfig(logger)
}

// @Description
func loadConfig(logger logger.Logger) *Config {
	cfg := &Config{
		RedisConfig: loadRedisConfig(logger),
	}
	return cfg
}

func getStringEnvOrDefault(key string, fallback string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return fallback
}

func getIntEnvOrDefault(key string, fallback int) int {
	if env := os.Getenv(key); env != "" {
		if i, err := strconv.Atoi(env); err == nil {
			return i
		}
	}
	return fallback
}

func getDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if env := os.Getenv(key); env != "" {
		if duration, err := time.ParseDuration(env); err != nil {
			return duration
		}
	}
	return fallback
}

type configKey string

var (
	rEDIS_PROTOCOL configKey = "REDIS_PROTOCOL"
	rEDIS_USERNAME configKey = "REDIS_USERNAME"
	rEDIS_PASSWORD configKey = "REDIS_PASSWORD"
	rEDIS_HOST     configKey = "REDIS_HOST"
	rEDIS_PORT     configKey = "REDIS_PORT"
)
