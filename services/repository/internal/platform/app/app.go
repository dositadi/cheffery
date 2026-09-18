package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/repository/internal/bcrypt/bcryptapp"
	"github.com/dositadi/cheffery/services/repository/internal/config"
	"github.com/dositadi/cheffery/services/repository/internal/platform/postgres"
	"github.com/dositadi/cheffery/services/repository/internal/user/userapp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userhttp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	cfg       *config.ApplicationConfig
	logger    logger.Logger
	server    *userhttp.Server
	dbPool    *pgxpool.Pool
	dbMetrics *postgres.Metrics
	canceFunc context.CancelFunc
}

func New() *App {
	app := &App{}

	app.logger = logger.New(os.Stdout)
	app.cfg = config.LoadAppConfig(app.logger)

	ctx, cancel := context.WithCancel(context.Background())
	app.canceFunc = cancel

	app.dbPool, app.dbMetrics = postgres.New(app.logger, app.cfg.PostgresEnv, app.cfg.RetryCfg).InitDB(ctx)

	validate := validator.New(validator.WithPrivateFieldValidation(), validator.WithRequiredStructEnabled(), validator.WithTagNameFuncBlankOmit())
	repo := userpostgres.New(app.logger, app.dbPool, app.cfg.RetryCfg)
	bcrypt := bcryptapp.New(app.logger)
	userUsecase := userapp.New(app.logger, repo, validate, bcrypt)
	dbHealth := postgres.NewHealthChecker(app.logger, app.dbPool, app.dbMetrics, app.cfg.RetryCfg.MaxWait)
	app.server = userhttp.New(app.logger, userUsecase, dbHealth)

	return app
}

func (a *App) StartServer() {
	scope := "app.startServer"
	reqID := "start-server"

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", a.cfg.ServerPort))
	if err != nil {
		a.logger.PrintFatal(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: fmt.Errorf("Failed to listen: %w", err).Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer(grpc.ConnectionTimeout(time.Duration(2 * time.Second)))

	//Register the service implementation to the protobuf
	repository.RegisterRepositoryServer(grpcServer, a.server)

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGTERM, syscall.SIGINT)
	chErr := make(chan error)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			chErr <- err
		}
	}()

	select {
	case e := <-chErr:
		a.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: e.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	case <-chSignal:
	}

	a.canceFunc()
	grpcServer.GracefulStop()
	a.dbPool.Close()
}
