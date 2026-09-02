package jwtserver

import (
	"context"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
)

type Executor interface {
	ExecuteGenerateTokenPair(ctx context.Context, arg jwtapp.ExecuteGenerateTokenPairInput) (jwtapp.ExecuteGenerateTokenPairOutput, error)
	ExecuteValidateAccessToken(ctx context.Context, reqId, tokenString string) (*jwtdomain.JWTAccessClaim, error)
	ExecuteValidateRefreshToken(ctx context.Context, reqId, refreshToken string) (*jwtdomain.JWTRefreshClaim, error)
	ExecuteLogout(ctx context.Context, arg jwtapp.ExecuteLogoutInput) (string, error)
	ExecuteRotateRefreshToken(ctx context.Context, arg jwtapp.ExecuteRotateRefreshTokenInput) (jwtapp.ExecuteRotateRefreshTokenOutput, error)
}
