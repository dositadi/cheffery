package config

import (
	"os"
	"strconv"
	"time"
)

// @Description This is the auth service config file,
// @Description It loads all the configurations this service will need to commence running
type Config struct {
	RedisConfig *RedisConfig
}

func LoadConfig() *Config {
	cfg := &Config{
		RedisConfig: loadRedisConfig(),
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
