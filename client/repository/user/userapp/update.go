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

type ExecuteUpdateInput struct {
	Name        string    `validate:"omitempty, gte=3"`
	Email       string    `validate:"omitempty, email"`
	Password    string    `validate:"omitempty, gte=8"`
	OldPassword string    `validate:"omitempty, gte=8"`
	ID          uuid.UUID `validate:"required, uuid"`
}

func (e ExecuteUpdateInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "Name":
					return customerror.WrapValidateErr(userdomain.ErrName, e.Error())
				case "Email":
					return customerror.WrapValidateErr(userdomain.ErrEmail, e.Error())
				case "Password":
					return customerror.WrapValidateErr(userdomain.ErrPassword, e.Error())
				case "ID":
					return customerror.WrapValidateErr(userdomain.ErrID, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(userdomain.ErrInternal, err.Error())
	}
	return nil
}

type ExecuteUpdateOutput struct {
	DataChanged bool
	Fields      []string
	Timestamp   time.Time
}

func (u *Usecase) ExecuteUpdate(ctx context.Context, arg ExecuteUpdateInput) (ExecuteUpdateOutput, error) {
	scope := "userapp.ExecuteUpdate"
	reqID := middleware.GetReqID(ctx)

	if err := arg.validate(u.validate); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		return ExecuteUpdateOutput{}, err
	}

	response, err := u.repoClient.UpdateUser(ctx, &repository.UpdateUserRequest{
		Name:        arg.Name,
		Email:       arg.Email,
		NewPassword: arg.Password,
		OldPassword: arg.OldPassword,
		Id:          arg.ID.String(),
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
			return ExecuteUpdateOutput{}, userdomain.ErrInternal
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrBadRequest, st.Message())
		case codes.Unauthenticated:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrUnauthorized, st.Message())
		case codes.AlreadyExists:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrEmailConflict, st.Message())
		case codes.NotFound:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrNotFound, st.Message())
		case codes.DeadlineExceeded:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrTimeout, st.Message())
		default:
			return ExecuteUpdateOutput{}, customerror.WrapValidateErr(userdomain.ErrInternal, st.Message())
		}
	}

	return ExecuteUpdateOutput{
		DataChanged: response.GetDatachanged(),
		Fields:      response.GetFields(),
		Timestamp:   response.GetTimestamp().AsTime(),
	}, nil
}
