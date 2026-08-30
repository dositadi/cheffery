package store

import (
	"fmt"

	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	logger logger.Logger
	client *redis.Client
}

func NewTokenStore(logger logger.Logger, client *redis.Client) *TokenStore {
	return &TokenStore{
		logger: logger,
		client: client,
	}
}

func tokenKey(tokenId string) string {
	return fmt.Sprintf("token:%s", tokenId)
}

func userKey(userId string) string {
	return fmt.Sprintf("userToken:%s", userId)
}

func tokenVersionKey(userId string) string {
	return fmt.Sprintf("token_version:%s", userId)
}
