package jwtserver

import (
	"context"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
)

type Executor interface {
	ExecuteGenerateAccessToken(ctx context.Context, reqId string, userId string) (string, error)
	ExecuteGenerateRefreshToken(ctx context.Context, reqId, userId string) (string, error)
	ExecuteValidateAccessToken(ctx context.Context, reqId, tokenString string) (*jwtdomain.JWTAccessClaim, error) 
	ExecuteValidateRefreshToken(ctx context.Context, reqId, refreshToken string) (*jwtdomain.JWTRefreshClaim, error)
}
