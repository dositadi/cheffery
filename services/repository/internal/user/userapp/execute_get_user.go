package userapp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ExecuteGetInput struct {
	ReqID  string    `validate:"required"`
	UserID uuid.UUID `validate:"required, eq=16"`
}

func (e ExecuteGetInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "ReqID":
					return userdomain.ErrReqID
				case "UserID":
					return userdomain.ErrID
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

func (u *Usecase) ExecuteGet(ctx context.Context, arg ExecuteGetInput) (*userdomain.User, error) {
	scope := "userapp.ExecuteGet"

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
		return nil, err
	}

	response, err := u.repo.GetUser(ctx, userpostgres.GetUserInput{
		ReqID:  arg.ReqID,
		UserID: arg.UserID,
	})
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
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
