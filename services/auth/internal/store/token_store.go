package store

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/dositadi/cheffery/services/shared/customerror"
)

func (t *TokenStore) IncrementTokenVersion(ctx context.Context, reqID, userId string) error {
	scope := "store.IncrementUserTokenVersion()"
	if err := t.client.Incr(ctx, tokenVersionKey(userId)).Err(); err != nil {
		t.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}
	return nil
}

func (t *TokenStore) GetTokenVersion(ctx context.Context, reqID, userID string) (int64, error) {
	scope := "store.GetTokenVersion"

	version, err := t.client.Get(ctx, tokenVersionKey(userID)).Int64()
	if err != nil {
		t.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return -1, err
	}

	return version, nil
}

func (t *TokenStore) StoreRefreshToken(ctx context.Context, reqID, tokenId, userId string, expiresAt time.Time) error {
	scope := "store.StoreRefreshToken"
	data := TokenData{
		UserId:    userId,
		TokenId:   tokenId,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	json, err := json.Marshal(data)
	if err != nil {
		t.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	ttl := time.Until(expiresAt)

	// Store the token
	if err := t.client.Set(ctx, refreshTokenKey(tokenId), json, ttl).Err(); err != nil {
		t.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	// Store token in User set
	if err := t.client.SAdd(ctx, userTokensSetKey(userId), tokenId).Err(); err != nil {
		t.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	return nil
}

func (t *TokenStore) GetRefreshToken(ctx context.Context, reqId, tokenId string) (TokenData, error) {
	scope := "store.GetRefreshToken"

	data, err := t.client.Get(ctx, refreshTokenKey(tokenId)).Bytes()
	if err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return TokenData{}, err
	}

	var out TokenData

	if err := json.Unmarshal(data, &out); err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return TokenData{}, err
	}

	return out, nil
}

func (t *TokenStore) RevokeRefreshToken(ctx context.Context, reqId, tokenId string) error {
	scope := "store.RevokeRefreshToken"

	tokenData, err := t.GetRefreshToken(ctx, reqId, tokenId)
	if err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	tokenData.Revoked = true

	json, err := json.Marshal(tokenData)
	if err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	ttl := time.Until(tokenData.ExpiresAt)

	if ttl <= 0 {
		return nil // token has expired already
	}

	if err := t.client.Set(ctx, refreshTokenKey(tokenId), json, ttl).Err(); err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	return nil
}

func (t *TokenStore) RevokeAllRefreshToken(ctx context.Context, reqId, userId string) error {
	scope := "store.RevokeAllRefreshToken"

	tokenIds, err := t.client.SMembers(ctx, userTokensSetKey(userId)).Result()
	if err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	revokeGroup := new(sync.WaitGroup)

	revokeGroup.Add(len(tokenIds))

	for _, token := range tokenIds {
		go func(token string) {
			defer revokeGroup.Done()
			if err := t.RevokeAllRefreshToken(ctx, reqId, userId); err != nil {
				t.logger.PrintError(err, reqId, customerror.InternalError{
					Inner:   err,
					Message: err.Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
			}
		}(token)
	}

	revokeGroup.Wait()

	return nil
}

func (t *TokenStore) BlacklistAccessToken(ctx context.Context, reqId, tokenId string, expiresAt time.Time) error {
	scope := "store.BlacklistAccessToken"
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil // Token has expired already, no need to blacklist it
	}

	if err := t.client.Set(ctx, blacklistTokenKey(tokenId), tokenId, ttl).Err(); err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}
	return nil
}

func (t *TokenStore) AccessTokenIsBlackListed(ctx context.Context, reqId, tokenId string) (bool, error) {
	scope := "store.AccessTokenIsBlackListed"
	exists, err := t.client.Exists(ctx, blacklistTokenKey(tokenId)).Result()
	if err != nil {
		t.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return false, err
	}

	return exists > 0, nil
}
