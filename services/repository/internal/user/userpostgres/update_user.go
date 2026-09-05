package userpostgres

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/store/sqlc"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type UpdateUserInput struct {
	Name         string
	Email        string
	PasswordHash []byte
	ID           uuid.UUID
	Version      int32
}

func (r *Repository) UpdateUser(ctx context.Context, reqID string, arg UpdateUserInput) error {
	scope := "userpostgres.UpdateUser"

	querier := sqlc.New(r.pgPool)
	wait := r.retryCfg.MinWait
	var err error

updateUser:
	for attempt := range r.retryCfg.MaxAttempt {
		err = querier.UpdateUser(ctx, sqlc.UpdateUserParams(arg))
		if err == nil {
			break updateUser
		}

		if !customerror.IsRetryableError(err) {
			break updateUser
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
				break updateUser
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

		var pgErr *pgconn.PgError
		if errors.As(err, pgErr) {
			switch pgErr.ConstraintName {
			case "idx_email":
				return ErrEmailConflict
			}
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ErrRequestTimeout
		}
		return ErrInternal
	}

	return nil
}
