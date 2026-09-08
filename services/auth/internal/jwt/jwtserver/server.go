package jwtserver

import (
	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/config"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Server struct {
	gen.UnimplementedIssuerServer
	port   Executor
	store  jwtdomain.TokenStore
	logger logger.Logger
	cfg    *config.JWTConfig
}

func New(port Executor, logger logger.Logger, cfg *config.JWTConfig, store jwtdomain.TokenStore) *Server {
	return &Server{
		port:   port,
		logger: logger,
		store:  store,
		cfg:    cfg,
	}
}
