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

type ExecuteDeleteInput struct {
	ID uuid.UUID `validate:"required, uuid"`
}

func (e ExecuteDeleteInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "ID":
					return fmt.Errorf("%w: %s", userdomain.ErrID, e.Error())
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

func (u *Usecase) ExecuteDelete(ctx context.Context, arg ExecuteDeleteInput) error {
	scope := "userapp.ExecuteDelete"
	reqID := middleware.GetReqID(ctx)

	if err := arg.validate(u.validate); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	user, err := u.ExecuteGetByID(ctx, ExecuteGetByIDInput{
		UserID: arg.ID,
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	if err := u.repo.DeleteUser(ctx, userpostgres.DeleteUserInput{
		ID:      arg.ID,
		Version: user.GetVersion(),
	}); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		if errors.Is(err, userpostgres.ErrRequestTimeout) {
			return userdomain.ErrTimeout
		}
		return userdomain.ErrInternal
	}
	return nil
}
