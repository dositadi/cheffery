package userapp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ExecuteCreateInput struct {
	Name     string `validate:"required, gte=3"`
	Email    string `validate:"required, email"`
	Password string `validate:"password, gte=8"`
	ReqID    string `validate:"omitempty"`
}

func (e ExecuteCreateInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {

		if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validationErrs {
				switch e.StructField() {
				case "Name":
					return fmt.Errorf("%w: %s", userdomain.ErrName, e.Error())
				case "Email":
					return fmt.Errorf("%w: %s", userdomain.ErrEmail, e.Error())
				case "Password":
					return fmt.Errorf("%w: %s", userdomain.ErrPassword, e.Error())
				}
			}
		}

		return userdomain.ErrInternal
	}

	return nil
}

type ExecuteCreateOutput struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (u *Usecase) ExecuteCreate(ctx context.Context, arg ExecuteCreateInput) (ExecuteCreateOutput, error) {
	scope := "userapp.ExecuteCreate"

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
		return ExecuteCreateOutput{}, err
	}

	hashedPassword, err := u.bcrypt.GenerateHash([]byte(arg.Password))
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteCreateOutput{}, userdomain.ErrInternal
	}

	response, err := u.repo.CreateUser(ctx, userpostgres.CreateUserInput{
		Name:         arg.Name,
		Email:        arg.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, userpostgres.ErrEmailConflict) {
			return ExecuteCreateOutput{}, userdomain.ErrEmailConflict
		}
		if errors.Is(err, userpostgres.ErrRequestTimeout) {
			return ExecuteCreateOutput{}, userdomain.ErrTimeout
		}
		return ExecuteCreateOutput{}, userdomain.ErrInternal
	}

	return ExecuteCreateOutput(response), nil
}
