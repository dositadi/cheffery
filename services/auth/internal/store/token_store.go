package store

import (
	"context"
	"encoding/json"
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
	scope := "store.StoreRefreshToken()"
	data := struct {
		UserId    string
		TokenId   string
		ExpiresAt time.Time
		Revoked   bool
	}{
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
