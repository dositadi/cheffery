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
	JWTConfig   *JWTConfig   // Jwt config
}

// @Summary new cofig type
// @Description creates new config type
// @Param logger
// @Produce Config
func New(logger logger.Logger) *Config {
	return loadConfig(logger)
}

// @Summary load configurations
// @Description load environment Configurations
// @Param logger
// @Produce Config
func loadConfig(logger logger.Logger) *Config {
	cfg := &Config{
		RedisConfig: loadRedisConfig(logger),
		JWTConfig:   loadJWTConfig(logger),
	}
	return cfg
}

// @Summary get string env variables
// @Description get string env variables or fallback to default
// @Param key string true "Variable key"
// @Param fallback string true "Default value"
// @Produce string
func getStringEnvOrDefault(key string, fallback string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return fallback
}

// @Summary get int env variables
// @Description get int env variables or fallback to default
// @Param key string true "Variable key"
// @Param fallback int true "Default value"
// @Produce int
func getIntEnvOrDefault(key string, fallback int) int {
	if env := os.Getenv(key); env != "" {
		if i, err := strconv.Atoi(env); err == nil {
			return i
		}
	}
	return fallback
}

// @Summary get duration env variables
// @Description get duration env variables or fallback to default
// @Param key string true "Variable key"
// @Param fallback time.Duration true "Default value"
// @Produce time.Duration
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
	rEDIS_PROTOCOL        configKey = "REDIS_SCHEME"
	rEDIS_USERNAME        configKey = "REDIS_USERNAME"
	rEDIS_PASSWORD        configKey = "REDIS_PASSWORD"
	rEDIS_HOST            configKey = "REDIS_HOST"
	rEDIS_PORT            configKey = "REDIS_PORT"
	rEDIS_MAX_ACTIVE_CONN configKey = "REDIS_MAX_ACTIVE_CONN"
	rEDIS_MAX_IDLE_CONN   configKey = "REDIS_MAX_IDLE_CONN"
)
