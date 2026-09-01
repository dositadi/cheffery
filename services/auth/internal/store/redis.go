package store

import (
	"fmt"
	"time"

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

type TokenData struct {
	UserId    string
	TokenId   string
	ExpiresAt time.Time
	Revoked   bool
}

func refreshTokenKey(tokenId string) string {
	return fmt.Sprintf("refresh_token:%s", tokenId)
}

func userTokensSetKey(userId string) string {
	return fmt.Sprintf("user_token_set:%s", userId)
}

func tokenVersionKey(userId string) string {
	return fmt.Sprintf("token_version:%s", userId)
}

func blacklistTokenKey(tokenId string) string {
	return fmt.Sprintf("blacklist_access_token:%s", tokenId)
}
