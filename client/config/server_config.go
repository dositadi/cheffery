package config

import "time"

type ServerConfig struct {
	Port           uint16
	Host           string
	MaxHeaderBytes uint16
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
}

func loadServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		Port:           getIntOrDefault("SERVER_PORT", 8080),
		Host:           getStringOrDefault("SERVER_HOST", "0.0.0.0"),
		MaxHeaderBytes: getIntOrDefault("SERVER_MAX_HEADER_BYTES", 13427),
		ReadTimeout:    getDurationOrDefault("SERVER_READ_TIMEOUT", time.Duration(3*time.Second)),
		WriteTimeout:   getDurationOrDefault("SERVER_WRITE_TIMEOUT", time.Duration(5*time.Second)),
		IdleTimeout:    getDurationOrDefault("SERVER_IDLE_TIMEOUT", time.Duration(10*time.Second)),
	}

	return cfg
}
