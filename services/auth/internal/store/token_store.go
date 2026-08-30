package store

import (
	"context"

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
