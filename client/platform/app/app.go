package app

import (
	"os"

	"github.com/dositadi/cheffery/client/config"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.AppConfig
	router     *chi.Mux
	logger     logger.Logger
	authClient auth.IssuerClient
	repoClient repository.RepositoryClient
	conns      []*grpc.ClientConn
}

func New() *App {
	app := &App{
		logger: logger.New(os.Stdout),
	}

	app.cfg = config.LoadAppConfig(app.logger)
	authConn := establishConn(app.logger, toAddr(app.cfg.Service.AuthHost, app.cfg.Service.AuthPort))
	app.authClient = auth.NewIssuerClient(authConn)

	repoConn := establishConn(app.logger, toAddr(app.cfg.Service.RepoHost, app.cfg.Service.RepoPort))
	app.repoClient = repository.NewRepositoryClient(repoConn)
	app.router = chi.NewRouter()

	app.conns = append(app.conns, authConn, repoConn)

	return app
}
