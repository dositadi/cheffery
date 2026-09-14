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

	cfg.PrepareConn = func(ctx context.Context, c *pgx.Conn) (bool, error) {
		p.logger.PrintInfo("db:prepare-conn", fmt.Sprintf("Preparing connection to %s:%v", c.Config().Host, c.Config().Port), map[string]string{
			"Context": scope,
		})
		return true, nil
	}

	cfg.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
		p.logger.PrintInfo("db:before-connect", fmt.Sprintf("Connecting to %s:%v", cc.Host, cc.Port), map[string]string{
			"Context": scope,
		})
		return nil
	}

	cfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		p.logger.PrintInfo("db:after-connect", fmt.Sprintf("Connected to %s:%v", c.Config().Host, c.Config().Port), map[string]string{
			"Context": scope,
		})
		return nil
	}

	cfg.AfterRelease = func(c *pgx.Conn) bool {
		p.logger.PrintInfo("db:after-release", fmt.Sprintf("Connection to %s:%v released", c.Config().Host, c.Config().Port), map[string]string{
			"Context": scope,
		})
		return true
	}

	return p.connect(ctx, cfg)
}

func (p *Postgres) connect(ctx context.Context, cfg *pgxpool.Config) *pgxpool.Pool {
	scope := "postgres.connect"
	wait := p.retryCfg.MinWait
	for attempt := range p.retryCfg.MaxAttempt {
		dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			return dbPool
		}
		p.logger.PrintError(err, "init-db", err.Error(), map[string]string{
			"Context": scope,
			"Attempt": strconv.Itoa(attempt),
		})

		select {
		case <-ctx.Done():
			p.logger.PrintFatal(nil, "init-db", "Failed to connect to the database", map[string]string{
				"Context": scope,
			})
			return nil
		case <-time.After(wait):
			wait *= 2
			if wait > p.retryCfg.MaxWait {
				wait = p.retryCfg.MaxWait
			}
		}

	}
	p.logger.PrintFatal(nil, "init-db", "Failed to connect to the database", map[string]string{
		"Context": scope,
	})
	panic("")
}
