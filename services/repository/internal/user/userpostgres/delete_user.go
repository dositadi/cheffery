package userpostgres

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/store/sqlc"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
)

type DeleteUserInput struct {
	ID      uuid.UUID
	Version int32
}

func (r *Repository) DeleteUser(ctx context.Context, reqID string, arg DeleteUserInput) error {
	scope := "userpostgres.DeleteUser"

	querier := sqlc.New(r.pgPool)
	wait := r.retryCfg.MinWait
	var err error

deleteUser:
	for attempt := range r.retryCfg.MaxAttempt {
		err = querier.DeleteUser(ctx, sqlc.DeleteUserParams(arg))
		if err == nil {
			break deleteUser
		}

		if !customerror.IsRetryableError(err) {
			break deleteUser
		}

		customerror.LogAttempt(r.logger, err, reqID, attempt, scope)

		if attempt < r.retryCfg.MaxAttempt {
			select {
			case <-ctx.Done():
				r.logger.PrintError(ctx.Err(), reqID, customerror.InternalError{
					Inner:   ctx.Err(),
					Message: ctx.Err().Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
				err = ctx.Err()
				break deleteUser
			case <-time.After(wait):
				wait *= 2
				if wait > r.retryCfg.MaxWait {
					wait = r.retryCfg.MaxWait
				}
			}
		}
	}
	if err != nil {
		r.logger.PrintError(ctx.Err(), reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ErrRequestTimeout
		}
		return ErrInternal
	}

	return nil
}
