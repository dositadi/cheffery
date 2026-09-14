package config

import (
	"os"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type PostgresEnv struct {
	Protocol          string
	Username          string
	Password          string
	Host              string
	Port              uint16
	Database          string
	ConnTimeout       time.Duration
	MaxConnIdleTime   time.Duration
	MaxConnLifetime   time.Duration
	MaxConns          uint16
	MinConns          uint16
	HealthCheckPeriod time.Duration
	MinIdleConn       uint16
}

func loadPosgresEnv(logger logger.Logger) *PostgresEnv {
	cfg := &PostgresEnv{
		Protocol:          getStringOrDefault("PG_PROTOCOL", "postgresql"),
		Username:          getStringOrDefault("PG_USERNAME", "postgres"),
		Host:              getStringOrDefault("PG_HOST", "db.wwkprijhodhmngvrmvji.supabase.co"),
		Port:              getIntOrDefault("PG_PORT", 5432),
		Database:          getStringOrDefault("PG_DATABASE", "postgres"),
		ConnTimeout:       getDurationOrDefault("PG_CONN_TIMEOUT", time.Duration(2*time.Second)),
		MaxConnIdleTime:   getDurationOrDefault("PG_MAX_CONN_IDLE_TIME", time.Duration(5*time.Minute)),
		MaxConnLifetime:   getDurationOrDefault("PG_MAX_CONN_LIFETIME", time.Duration(30*time.Minute)),
		MaxConns:          getIntOrDefault("PG_MAX_CONNS", 30),
		MinConns:          getIntOrDefault("PG_MIN_CONNS", 10),
		HealthCheckPeriod: getDurationOrDefault("HEALTH_CHECK_PERIOD", time.Duration(2*time.Second)),
		MinIdleConn:       getIntOrDefault("MIN_IDLE_CONN", 20),
	}

	password := os.Getenv("PG_PASSWORD")
	if password == "" {
		logger.PrintFatal(nil, "loadPostgresEnv", "PG_PASSWORD environment is required", map[string]string{
			"Context": "config.loadPosgresEnv",
		})
	}

	cfg.Password = password

	return cfg
}
