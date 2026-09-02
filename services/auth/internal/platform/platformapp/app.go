package platformapp

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/config"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtserver"
	"github.com/dositadi/cheffery/services/auth/internal/platform/rdb/rdbdomain"
	"github.com/dositadi/cheffery/services/auth/internal/store"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/services/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *config.Config
	logger     logger.Logger
	redis      *rdbdomain.Redis
	tokenStore *store.TokenStore
	server     *jwtserver.Server
	executor   *jwtapp.Usecase
}

func New() *App {
	app := &App{}

	app.logger = logger.New(os.Stdout)
	app.cfg = config.New(app.logger)
	app.redis = rdbdomain.New(app.logger, app.cfg.RedisConfig)
	app.tokenStore = store.NewTokenStore(app.logger, app.redis.Client)
	app.executor = jwtapp.New(app.logger, app.cfg.JWTConfig, app.tokenStore)
	app.server = jwtserver.New(app.executor, app.logger, app.cfg.JWTConfig, app.tokenStore)

	return app
}

func (a *App) startServer() {
	scope := "platformapp.startServer"
	reqId := "start-server"
	// Listen to port
	listener, err := net.Listen("tcp", ":"+a.cfg.ServerPort)
	if err != nil {
		a.logger.PrintFatal(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: fmt.Errorf("failed to listen: %w", err).Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer(grpc.ConnectionTimeout(time.Duration(2 * time.Second)))

	// Register the service implementation to the protobuf
	gen.RegisterIssuerServer(grpcServer, a.server)

	// Enable reflection for debugging with tools like grpcurl.
	reflection.Register(grpcServer)

	chSignal := make(chan os.Signal, 1)
	chErr := make(chan error)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			chErr <- err
		}
	}()

	select {
	case err = <-chErr:
		a.logger.PrintFatal(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: fmt.Errorf("failed to listen: %w", err).Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	case <-chSignal:
	}

	a.logger.PrintInfo(reqId, "starting auth gRPC server on: :"+a.cfg.ServerPort, map[string]string{
		"Context": scope,
	})

	<-chSignal

	grpcServer.GracefulStop()
	a.redis.Client.Close()
}

func (a *App) Run() {
	a.startServer()
}
