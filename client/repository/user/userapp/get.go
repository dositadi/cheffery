package userapp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExecuteGetByEmailInput struct {
	Email string `validate:"required, email"`
}

func (e ExecuteGetByEmailInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {

		if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validationErrs {
				switch e.StructField() {
				case "Email":
					return customerror.WrapValidateErr(userdomain.ErrEmail, e.Error())
				}
			}
		}

		return customerror.WrapValidateErr(userdomain.ErrInternal, err.Error())
	}

	return nil
}

func (u *Usecase) ExecuteGetByEmail(ctx context.Context, arg ExecuteGetByEmailInput) (*userdomain.User, error) {
	scope := "userapp.ExecuteGetByEmail"
	reqID := middleware.GetReqID(ctx)

	if err := arg.validate(u.validate); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, err
	}

	response, err := u.repoClient.GetUserByEmail(ctx, &repository.GetUserByEmailRequest{
		Email: arg.Email,
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
			return nil, userdomain.ErrInternal
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, customerror.WrapValidateErr(userdomain.ErrNotFound, st.Message())
		case codes.DeadlineExceeded:
			return nil, customerror.WrapValidateErr(userdomain.ErrTimeout, st.Message())
		default:
			return nil, customerror.WrapValidateErr(userdomain.ErrInternal, st.Message())
		}
	}

	id, err := uuid.Parse(response.GetId())
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, userdomain.ErrInternal
	}

	user, err := userdomain.New(id, response.GetName(), response.GetEmail(), response.GetPasswordHash(), response.GetVersion(), response.GetCreatedAt().AsTime(), response.GetUpdatedAt().AsTime())
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, userdomain.ErrInternal
	}

	return user, nil
}

func (u *Usecase) ExecuteGetByID(ctx context.Context, id uuid.UUID) (*userdomain.User, error) {
	scope := "userapp.ExecuteGetByID"
	reqID := middleware.GetReqID(ctx)

	response, err := u.repoClient.GetUserByID(ctx, &repository.GetUserByIDRequest{
		UserID: id.String(),
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
			return nil, userdomain.ErrInternal
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, customerror.WrapValidateErr(userdomain.ErrNotFound, st.Message())
		case codes.DeadlineExceeded:
			return nil, customerror.WrapValidateErr(userdomain.ErrTimeout, st.Message())
		default:
			return nil, customerror.WrapValidateErr(userdomain.ErrInternal, st.Message())
		}
	}

	user, err := userdomain.New(id, response.GetName(), response.GetEmail(), nil, response.GetVersion(), response.GetCreatedAt().AsTime(), response.GetUpdatedAt().AsTime())
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, userdomain.ErrInternal
	}

	return user, nil
}
