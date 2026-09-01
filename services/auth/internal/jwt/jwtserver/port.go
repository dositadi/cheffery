package jwtserver

import (
	"context"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
)

type Executor interface {
	ExecuteGenerateAccessToken(ctx context.Context, reqId string, arg jwtapp.ExecuteGenerateAccessTokenInput) (string, error)
	ExecuteGenerateRefreshToken(ctx context.Context, reqId, userId string) (string, error)
	ExecuteValidateAccessToken(ctx context.Context, tokenString string) (*jwtdomain.JWTAccessClaim, error)
	ExecuteValidateRefreshToken(ctx context.Context, refreshToken string) (*jwtdomain.JWTRefreshClaim, error)
}
