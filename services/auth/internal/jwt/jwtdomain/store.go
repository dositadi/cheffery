package jwtdomain

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/store"
)

type TokenStore interface {
	IncrementTokenVersion(ctx context.Context, reqID, userId string) error
	GetTokenVersion(ctx context.Context, reqID, userID string) (int64, error)
	StoreRefreshToken(ctx context.Context, reqID, tokenId, userId string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, reqId, tokenId string) (store.TokenData, error)
	RevokeRefreshToken(ctx context.Context, reqId, tokenId string) error
	RevokeAllRefreshToken(ctx context.Context, reqId, userId string) error
	BlacklistAccessToken(ctx context.Context, reqId, tokenId string, expiresAt time.Time) error
	AccessTokenIsBlackListed(ctx context.Context, reqId, tokenId string) (bool, error)
}
