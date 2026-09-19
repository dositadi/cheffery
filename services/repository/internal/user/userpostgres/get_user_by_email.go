package userpostgres

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/store/sqlc"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

type GetUserByEmailInput struct {
	Email string
}

func (r *Repository) GetUserByEmail(ctx context.Context, arg GetUserByEmailInput) (GetUserOutput, error) {
	scope := "userpostgres.GetUser"
	reqID := middleware.GetReqID(ctx)

	querier := sqlc.New(r.pgPool)

	wait := r.retryCfg.MinWait
	var err error
	var response sqlc.User

getUser:
	for attempt := range r.retryCfg.MaxAttempt {
		response, err = querier.GetUserByEmail(ctx, arg.Email)
		if err == nil {
			break getUser
		}

		if !customerror.IsRetryableError(err) {
			break getUser
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
				break getUser
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

		if errors.Is(err, pgx.ErrNoRows) {
			return GetUserOutput{}, ErrNotFound
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return GetUserOutput{}, ErrRequestTimeout
		}
		return GetUserOutput{}, ErrInternal
	}

	var deletedAt *time.Time
	if response.DeletedAt.Valid {
		deletedAt = &response.DeletedAt.Time
	}

	return GetUserOutput{
		ID:           response.ID,
		Name:         response.Name,
		Email:        response.Email,
		PasswordHash: response.PasswordHash,
		Createdat:    response.CreatedAt,
		Updatedat:    response.UpdatedAt,
		Version:      response.Version,
		Deletedat:    deletedAt,
	}, nil
}
