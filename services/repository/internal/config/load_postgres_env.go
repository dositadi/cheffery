package config

import (
	"os"

	"github.com/dositadi/cheffery/services/shared/logger"
)

type PostgresEnv struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func loadPosgresEnv(logger logger.Logger) *PostgresEnv {
	cfg := &PostgresEnv{
		Protocol: getStringOrDefault("PG_PROTOCOL", "postgresql"),
		Username: getStringOrDefault("PG_USERNAME", "postgres"),
		Host:     getStringOrDefault("PG_HOST", "db.wwkprijhodhmngvrmvji.supabase.co"),
		Port:     getStringOrDefault("PG_PORT", "5432"),
		Database: getStringOrDefault("PG_DATABASE", "postgres"),
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
