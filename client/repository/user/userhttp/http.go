package userhttp

import (
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Handler struct {
	logger     logger.Logger
	port       UsecasePort
	authClient auth.IssuerClient
}

func New(logger logger.Logger, port UsecasePort, authClient auth.IssuerClient) *Handler {
	return &Handler{
		port:       port,
		logger:     logger,
		authClient: authClient,
	}
}
