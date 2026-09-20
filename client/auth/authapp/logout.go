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

func (u *Usecase) Logout(ctx context.Context, accessToken sharedkernel.JWT) error {
	scope := "authapp.Logout"
	reqID := middleware.GetReqID(ctx)

	_, err := u.authService.Logout(ctx, &auth.LogoutRequest{
		AccessToken: accessToken.String(),
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
			return authdomain.ErrInternal
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return authdomain.ErrUnauthorized
		case codes.Internal:
			return authdomain.ErrInternal
		case codes.DeadlineExceeded:
			return authdomain.ErrTimeout
		}
		return authdomain.ErrInternal
	}

	return nil
}
