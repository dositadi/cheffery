package jwtserver

import (
	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/config"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Server struct {
	gen.UnimplementedIssuerServer
	port   Executor
	logger logger.Logger
	cfg    *config.JWTConfig
}

func New(port Executor, logger logger.Logger, cfg *config.JWTConfig) *Server {
	return &Server{
		port:   port,
		logger: logger,
		cfg:    cfg,
	}
}
