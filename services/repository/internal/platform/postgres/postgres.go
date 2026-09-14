package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/config"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	logger   logger.Logger
	pgCfg    *config.PostgresEnv
	retryCfg *config.RetryConfig
}

func New(logger logger.Logger, cfg *config.PostgresEnv, retryCfg *config.RetryConfig) *Postgres {
	return &Postgres{
		logger:   logger,
		pgCfg:    cfg,
		retryCfg: retryCfg,
	}
}

func (p *Postgres) InitDB(ctx context.Context) *pgxpool.Pool {
	scope := "postgres.InitDB"
	connConfig := &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:           p.pgCfg.Host,
			Port:           p.pgCfg.Port,
			Database:       p.pgCfg.Database,
			User:           p.pgCfg.Username,
			Password:       p.pgCfg.Password,
			ConnectTimeout: p.pgCfg.ConnTimeout,
		},
	}

	cfg := &pgxpool.Config{
		ConnConfig:        connConfig,
		MaxConnLifetime:   p.pgCfg.MaxConnLifetime,
		MaxConnIdleTime:   p.pgCfg.MaxConnIdleTime,
		MaxConns:          int32(p.pgCfg.MaxConns),
		MinConns:          int32(p.pgCfg.MinConns),
		HealthCheckPeriod: p.pgCfg.HealthCheckPeriod,
		MinIdleConns:      int32(p.pgCfg.MinIdleConn),
	}

	cfg.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
		p.logger.PrintInfo("db:before-connect", fmt.Sprintf("Connecting to %s:%v", cfg.ConnConfig.Host, cfg.ConnConfig.Port), map[string]string{
			"Context": scope,
		})
		return nil
	}

	return p.connect(ctx, cfg)
}

func (p *Postgres) connect(ctx context.Context, cfg *pgxpool.Config) *pgxpool.Pool {
	for attempt := range p.retryCfg.MaxAttempt {
		dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			return dbPool
		}
		p.logger.PrintError(err, "init-db", err.Error(), map[string]string{
			"Context": "postgres.connect",
			"Attempt": strconv.Itoa(attempt),
		})
		time.Sleep(p.retryCfg.MaxWait)
	}
	p.logger.PrintFatal(nil, "init-db", "Failed to connect to the database", map[string]string{
		"Context": "postgres.connect",
	})
	panic("")
}
