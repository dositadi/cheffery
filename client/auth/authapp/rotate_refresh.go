package authapp

import (
	"context"

	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/sharedkernel"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *Usecase) RotateRefresh(ctx context.Context, refreshToken sharedkernel.JWT) (*authdomain.TokenPair, error) {
	scope := "authapp.RotateRefresh"
	reqID := middleware.GetReqID(ctx)

	tokenPair, err := u.authService.RotateRefreshToken(ctx, &auth.RotateRefreshTokenRequest{
		RefreshToken: refreshToken.String(),
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		st, ok := status.FromError(err)
		if !ok {
			return nil, authdomain.ErrInternal
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return nil, authdomain.ErrUnauthorized
		case codes.DeadlineExceeded:
			return nil, authdomain.ErrTimeout
		case codes.Internal:
			return nil, authdomain.ErrInternal
		}
		return nil, authdomain.ErrInternal
	}

	out, err := authdomain.NewTokenPair(tokenPair.AccessToken, tokenPair.RefreshToken, tokenPair.ExpiresAt.AsTime(), u.validate)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, authdomain.ErrInternal
	}
	return out, nil
}
