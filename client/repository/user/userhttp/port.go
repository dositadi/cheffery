package userhttp

import (
	"context"

	"github.com/dositadi/cheffery/client/repository/user/userapp"
	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/google/uuid"
)

type UsecasePort interface {
	ExecuteCreate(ctx context.Context, arg userapp.ExecuteCreateInput) (userapp.ExecuteCreateOutput, error)
	ExecuteDelete(ctx context.Context, id uuid.UUID) error
	GetUserByID(ctx context.Context, req *repository.GetUserByIDRequest) (*repository.GetUserByIDResponse, error)
	ExecuteGetByEmail(ctx context.Context, arg userapp.ExecuteGetByEmailInput) (*userdomain.User, error)
	ExecuteGetByID(ctx context.Context, id uuid.UUID) (*userdomain.User, error)
	ExecuteUpdate(ctx context.Context, arg userapp.ExecuteUpdateInput) (userapp.ExecuteUpdateOutput, error)
}
