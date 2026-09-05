package userpostgres

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/store/sqlc"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GetUserInput struct {
	ReqID  string
	UserID uuid.UUID
}

type GetUserOutput struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash []byte
	Createdat    time.Time
	Updatedat    time.Time
	Version      int32
	Deletedat    *time.Time
}

func (r *Repository) GetUser(ctx context.Context, arg GetUserInput) (GetUserOutput, error) {
	scope := "userpostgres.GetUser"

	querier := sqlc.New(r.pgPool)

	wait := r.retryCfg.MinWait
	var err error
	var response sqlc.User

getUser:
	for attempt := range r.retryCfg.MaxAttempt {
		response, err = querier.GetUser(ctx, arg.UserID)
		if err == nil {
			break getUser
		}

		if !customerror.IsRetryableError(err) {
			break getUser
		}

		customerror.LogAttempt(r.logger, err, arg.ReqID, attempt, scope)

		if attempt < r.retryCfg.MaxAttempt {
			select {
			case <-ctx.Done():
				r.logger.PrintError(ctx.Err(), arg.ReqID, customerror.InternalError{
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
		r.logger.PrintError(ctx.Err(), arg.ReqID, customerror.InternalError{
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
	if response.Deletedat.Valid {
		deletedAt = &response.Deletedat.Time
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
