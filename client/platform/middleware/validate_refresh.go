package middleware

import (
	"context"
	"net/http"

	"github.com/dositadi/cheffery/globalutil/httphelper"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/sharedkernel"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *PlatformMiddleware) ValidateRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scope := "middleware.ValidateRefresh"
		reqID := middleware.GetReqID(r.Context())

		refreshToken := httphelper.GetRefreshToken(r)

		response, err := m.authClient.ValidateRefreshToken(r.Context(), &auth.ValidateRefreshTokenRequest{
			RefreshToken: refreshToken,
		})
		if err != nil {
			m.logger.PrintInfo(reqID, customerror.InternalError{
				Inner:   err,
				Message: err.Error(),
				Misc:    nil,
			}.Error(), map[string]string{
				"Context": scope,
			})

			err, ok := status.FromError(err)
			if !ok {
				customerror.InternalResponse(w, r, ErrInternal.Error())
				return
			}

			switch err.Code() {
			case codes.Unauthenticated:
				customerror.UnauthorizedResponse(w, r, ErrUnauthorized.Error())
				return
			case codes.DeadlineExceeded:
				customerror.RequestTimeoutResponse(w, r, ErrTimeout.Error())
				return
			case codes.Internal:
				customerror.InternalResponse(w, r, ErrInternal.Error())
				return
			}

			customerror.InternalResponse(w, r, ErrInternal.Error())
		}

		claim := sharedkernel.UserClaim{
			UserID: response.GetUserID(),
		}

		ctx := context.WithValue(r.Context(), sharedkernel.Claim, claim)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
