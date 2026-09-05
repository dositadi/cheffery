package userdomain

import (
	"context"

	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
)

type Repository interface {
	GetUser(ctx context.Context, arg userpostgres.GetUserInput) (userpostgres.GetUserOutput, error)
	CreateUser(ctx context.Context, reqId string, in userpostgres.CreateUserInput) (userpostgres.CreateUserOutput, error)
	UpdateUser(ctx context.Context, reqID string, arg userpostgres.UpdateUserInput) error
}
