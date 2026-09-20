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
	"github.com/google/uuid"
)

type ExecuteGetByIDInput struct {
	UserID uuid.UUID `validate:"required, uuid"`
}

func (e ExecuteGetByIDInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
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

func (u *Usecase) ExecuteGetByID(ctx context.Context, arg ExecuteGetByIDInput) (*userdomain.User, error) {
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

	response, err := u.repo.GetUserByID(ctx, userpostgres.GetUserByIDInput{
		UserID: arg.UserID,
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
