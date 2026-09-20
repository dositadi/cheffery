package middleware

import (
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type PlatformMiddleware struct {
	logger     logger.Logger
	authClient auth.IssuerClient
}

func NewPlatformMW(logger logger.Logger, authClient auth.IssuerClient) *PlatformMiddleware {
	return &PlatformMiddleware{
		logger:     logger,
		authClient: authClient,
	}
}
