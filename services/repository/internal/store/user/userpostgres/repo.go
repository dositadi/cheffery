package userpostgres

import (
	"github.com/dositadi/cheffery/services/repository/internal/config"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger   logger.Logger
	pgPool   *pgxpool.Pool
	retryCfg *config.RetryConfig
}

func New(logger logger.Logger, pgPool *pgxpool.Pool, retryCfg *config.RetryConfig) *Repository {
	return &Repository{
		logger:   logger,
		pgPool:   pgPool,
		retryCfg: retryCfg,
	}
}
