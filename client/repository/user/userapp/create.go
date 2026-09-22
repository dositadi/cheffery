package userapp

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExecuteCreateInput struct {
	Name     string `validate:"required, gte=3"`
	Email    string `validate:"required, email"`
	Password string `validate:"password, gte=8"`
}

func (e ExecuteCreateInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {

		if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validationErrs {
				switch e.StructField() {
				case "Name":
					return customerror.WrapValidateErr(userdomain.ErrName, e.Error())
				case "Email":
					return customerror.WrapValidateErr(userdomain.ErrEmail, e.Error())
				case "Password":
					return customerror.WrapValidateErr(userdomain.ErrPassword, e.Error())
				}
			}
		}

		return customerror.WrapValidateErr(userdomain.ErrInternal, err.Error())
	}

	return nil
}

type ExecuteCreateOutput struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (u *Usecase) ExecuteCreate(ctx context.Context, arg ExecuteCreateInput) (ExecuteCreateOutput, error) {
	scope := "userapp.ExecuteCreate"
	reqID := middleware.GetReqID(ctx)

	if err := arg.validate(u.validate); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteCreateOutput{}, err
	}

	response, err := u.repoClient.CreateUser(ctx, &repository.CreateUserRequest{
		Name:     arg.Name,
		Email:    arg.Email,
		Password: arg.Password,
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		err, ok := status.FromError(err)
		if !ok {
			return ExecuteCreateOutput{}, userdomain.ErrInternal
		}

		switch err.Code() {
		case codes.InvalidArgument:
			return ExecuteCreateOutput{}, customerror.WrapValidateErr(userdomain.ErrBadRequest, err.String())
		case codes.Unauthenticated:
			return ExecuteCreateOutput{}, userdomain.ErrUnauthorized
		case codes.AlreadyExists:
			return ExecuteCreateOutput{}, userdomain.ErrEmailConflict
		case codes.DeadlineExceeded:
			return ExecuteCreateOutput{}, userdomain.ErrTimeout
		case codes.Internal:
			return ExecuteCreateOutput{}, userdomain.ErrInternal
		}
		return ExecuteCreateOutput{}, userdomain.ErrInternal
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
		return ExecuteCreateOutput{}, userdomain.ErrInternal
	}

	u.logger.PrintInfo(reqID, "Created a new user", map[string]string{
		"Context": scope,
		"ID":      id.String(),
	})

	return ExecuteCreateOutput{
		ID:        id,
		CreatedAt: response.GetCreatedAt().AsTime(),
	}, nil
}
