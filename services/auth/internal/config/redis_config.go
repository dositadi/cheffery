package config

import (
	"log"
	"os"
)

type RedisConfig struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     int
}

func loadRedisConfig() *RedisConfig {
	cfg := &RedisConfig{
		Protocol: getStringEnvOrDefault(string(rEDIS_PROTOCOL), "redis"),
		Username: getStringEnvOrDefault(string(rEDIS_USERNAME), "default"),
		Host:     getStringEnvOrDefault(string(rEDIS_HOST), "machine-adequate-crayon-70404.db.redis.io"),
		Port:     getIntEnvOrDefault(string(rEDIS_PORT), 18608),
	}

	cfg.Password = os.Getenv(string(rEDIS_PASSWORD))
	if cfg.Password == "" {
		// implement logger
		log.Fatalf("%s environment variable is required", rEDIS_PASSWORD)
	}

	return cfg
}
