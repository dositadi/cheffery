package userapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

type ExecuteGetByEmailInput struct {
	Email string `validate:"required, email"`
}

func (e ExecuteGetByEmailInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "UserID":
					return fmt.Errorf("%w: %s", userdomain.ErrID, e.Error())
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

func (u *Usecase) ExecuteGetByEmail(ctx context.Context, arg ExecuteGetByEmailInput) (*userdomain.User, error) {
	scope := "userapp.ExecuteGet"
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

	response, err := u.repo.GetUserByEmail(ctx, userpostgres.GetUserByEmailInput{
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
		if errors.Is(err, userpostgres.ErrNotFound) {
			return nil, userdomain.ErrNotFound
		}
		if errors.Is(err, userpostgres.ErrRequestTimeout) {
			return nil, userdomain.ErrTimeout
		}
		return nil, userdomain.ErrInternal
	}

	user, err := userdomain.New(response.ID, response.Name, response.Email, response.PasswordHash, response.Version, response.Createdat, response.Updatedat)
	if err != nil {
		return nil, userdomain.ErrInternal
	}

	return user, nil
}
