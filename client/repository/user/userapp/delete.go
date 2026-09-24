package userapp

import (
	"context"

	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *Usecase) ExecuteDelete(ctx context.Context, id uuid.UUID) error {
	scope := "userapp.ExecuteDelete"
	reqID := middleware.GetReqID(ctx)

	if _, err := u.repoClient.DeleteUser(ctx, &repository.DeleteUserRequest{
		Id: id.String(),
	}); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		st, ok := status.FromError(err)
		if !ok {
			return userdomain.ErrInternal
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return customerror.WrapValidateErr(userdomain.ErrUnauthorized, st.Message())
		case codes.NotFound:
			return customerror.WrapValidateErr(userdomain.ErrNotFound, st.Message())
		case codes.DeadlineExceeded:
			return customerror.WrapValidateErr(userdomain.ErrTimeout, st.Message())
		default:
			return customerror.WrapValidateErr(userdomain.ErrInternal, st.Message())
		}
	}

	return nil
}
