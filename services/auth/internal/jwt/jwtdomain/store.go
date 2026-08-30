package jwtdomain

import "context"

type TokenStore interface {
	IncrementTokenVersion(ctx context.Context, reqID, userId string) error
	GetTokenVersion(ctx context.Context, reqID, userID string) (int64, error)
}
