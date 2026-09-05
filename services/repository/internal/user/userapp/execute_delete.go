package userapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ExecuteDeleteInput struct {
	ID    uuid.UUID `validate:"required, eq=16"`
	ReqID string    `validate:"required"`
}

func (e ExecuteDeleteInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "ID":
					return fmt.Errorf("%w: %s", userdomain.ErrID, e.Error())
				case "ReqID":
					return fmt.Errorf("%w: %s", userdomain.ErrReqID, e.Error())
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

func (u *Usecase) ExecuteDelete(ctx context.Context, arg ExecuteDeleteInput) error {
	scope := "userapp.ExecuteDelete"

	if err := arg.validate(u.validate); err != nil {
		if arg.ReqID == "" {
			arg.ReqID = "userapp.ExecuteCreate-Request"
		}
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	user, err := u.ExecuteGet(ctx, ExecuteGetInput{
		ReqID:  arg.ReqID,
		UserID: arg.ID,
	})
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return err
	}

	if err := u.repo.DeleteUser(ctx, arg.ReqID, userpostgres.DeleteUserInput{
		ID:      arg.ID,
		Version: user.GetVersion(),
	}); err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
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
