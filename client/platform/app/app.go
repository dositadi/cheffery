package app

import (
	"os"

	"github.com/dositadi/cheffery/client/config"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	cfg        *config.AppConfig
	router     *chi.Mux
	logger     logger.Logger
	authClient auth.IssuerClient
	repoClient repository.RepositoryClient
}

func New() *App {
	app := &App{
		cfg:    config.LoadAppConfig(),
		logger: logger.New(os.Stdout),
	}

	app.authClient = auth.NewIssuerClient(establishConn(app.logger, ""))
	app.repoClient = repository.NewRepositoryClient(establishConn(app.logger, ""))
	app.router = chi.NewRouter()

	return app
}
