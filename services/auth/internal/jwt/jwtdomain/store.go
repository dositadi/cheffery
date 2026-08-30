package jwtdomain

import (
	"context"
	"time"
)

type TokenStore interface {
	IncrementTokenVersion(ctx context.Context, reqID, userId string) error
	GetTokenVersion(ctx context.Context, reqID, userID string) (int64, error)
	StoreRefreshToken(ctx context.Context, reqID, tokenId, userId string, expiresAt time.Time) error
}
