package jwtapp

import (
	"github.com/dositadi/cheffery/services/auth/internal/config"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/logger"
)

type Usecase struct {
	logger logger.Logger
	cfg    config.JWTConfig
	store  jwtdomain.TokenStore
}
