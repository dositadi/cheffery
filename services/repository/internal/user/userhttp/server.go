package userhttp

import (
	"context"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/repository/internal/platform/postgres"
	"github.com/dositadi/cheffery/services/repository/internal/user/userapp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Executors interface {
	ExecuteCreate(ctx context.Context, arg userapp.ExecuteCreateInput) (userapp.ExecuteCreateOutput, error)
	ExecuteDelete(ctx context.Context, arg userapp.ExecuteDeleteInput) error
	ExecuteGetByID(ctx context.Context, arg userapp.ExecuteGetByIDInput) (*userdomain.User, error)
	ExecuteUpdate(ctx context.Context, arg userapp.ExecuteUpdateInput) (userapp.ExecuteUpdateOutput, error)
	ExecuteGetByEmail(ctx context.Context, arg userapp.ExecuteGetByEmailInput) (*userdomain.User, error)
}

type DBHealthChecker interface {
	Check(ctx context.Context) *postgres.HealthStatus
}

type Server struct {
	repository.UnimplementedRepositoryServer
	logger        logger.Logger
	executor      Executors
	healthChecker DBHealthChecker
}

func New(logger logger.Logger, executor Executors, healthChecker DBHealthChecker) *Server {
	return &Server{
		logger:        logger,
		executor:      executor,
		healthChecker: healthChecker,
	}
}
