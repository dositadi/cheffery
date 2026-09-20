package authhttp

import (
	"context"

	"github.com/dositadi/cheffery/client/auth/authapp"
	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/sharedkernel"
)

type Usecase interface {
	Login(ctx context.Context, arg authapp.LoginInput) (*authdomain.TokenPair, error)
	Logout(ctx context.Context, accessToken sharedkernel.JWT) error
	RotateRefresh(ctx context.Context, refreshToken sharedkernel.JWT) (*authdomain.TokenPair, error)
}
