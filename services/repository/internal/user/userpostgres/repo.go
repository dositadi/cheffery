package userpostgres

import (
	"errors"

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

var (
	ErrEmailConflict  = errors.New("this email exists already")
	ErrRequestTimeout = errors.New("request timed out")
	ErrInternal       = errors.New("something wrong happened")
	ErrNoRows         = errors.New("no user found")
	ErrNotFound       = errors.New("user not found")
)
