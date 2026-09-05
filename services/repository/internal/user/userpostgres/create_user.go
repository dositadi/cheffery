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

type CreateUserInput struct {
	Name         string
	Email        string
	PasswordHash []byte
}

type CreateUserOutput struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (r *Repository) CreateUser(ctx context.Context, reqId string, in CreateUserInput) (CreateUserOutput, error) {
	scope := "userpostgres.CreateUser"

	querier := sqlc.New(r.pgPool)

	wait := r.retryCfg.MinWait
	var err error
	var response sqlc.CreateUserRow

createUser:
	for attempt := range r.retryCfg.MaxAttempt {
		response, err = querier.CreateUser(ctx, sqlc.CreateUserParams(in))
		if err == nil {
			break createUser
		}

		if !customerror.IsRetryableError(err) {
			break createUser
		}

		customerror.LogAttempt(r.logger, err, reqId, attempt, scope)

		if attempt < r.retryCfg.MaxAttempt {
			select {
			case <-ctx.Done():
				r.logger.PrintError(ctx.Err(), reqId, customerror.InternalError{
					Inner:   ctx.Err(),
					Message: ctx.Err().Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
				err = ctx.Err()
				break createUser
			case <-time.After(wait):
				wait *= 2
				if wait > r.retryCfg.MaxWait {
					wait = r.retryCfg.MaxWait
				}
			}
		}
	}
	if err != nil {
		r.logger.PrintError(ctx.Err(), reqId, customerror.InternalError{
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
				return CreateUserOutput{}, ErrEmailConflict
			}
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return CreateUserOutput{}, ErrRequestTimeout
		}
		return CreateUserOutput{}, ErrInternal
	}
	return CreateUserOutput{
		ID:        response.ID,
		CreatedAt: response.CreatedAt,
	}, nil
}
