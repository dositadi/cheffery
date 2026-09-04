package userpostgres

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/store/sqlc"
	"github.com/google/uuid"
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

func (r *Repository) CreateUser(ctx context.Context, reqId, in CreateUserInput) (CreateUserOutput, error) {
	scope := "userpostgres.CreateUser"

	querier := sqlc.New(r.pgPool)

	wait := r.retryCfg.MinWait
	var err error
	var response sqlc.CreateUserRow

createUser:
	for attempt := range r.retryCfg.MaxAttempt {
		response, err = querier.CreateUser(ctx, sqlc.CreateUserParams(in))
		
	}
}
