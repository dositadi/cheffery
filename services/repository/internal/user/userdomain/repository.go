package userdomain

import (
	"context"

	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
)

type Repository interface {
	GetUserByID(ctx context.Context, arg userpostgres.GetUserByIDInput) (userpostgres.GetUserOutput, error)
	GetUserByEmail(ctx context.Context, arg userpostgres.GetUserByEmailInput) (userpostgres.GetUserOutput, error)
	CreateUser(ctx context.Context, in userpostgres.CreateUserInput) (userpostgres.CreateUserOutput, error)
	UpdateUser(ctx context.Context, arg userpostgres.UpdateUserInput) error
	DeleteUser(ctx context.Context, arg userpostgres.DeleteUserInput) error
}
