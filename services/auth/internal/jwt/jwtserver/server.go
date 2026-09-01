package jwtserver

import (
	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Server struct {
	gen.UnimplementedIssuerServer
	port   Executor
	logger logger.Logger
}

func New(port Executor, logger logger.Logger) *Server {
	return &Server{
		port:   port,
		logger: logger,
	}
}


