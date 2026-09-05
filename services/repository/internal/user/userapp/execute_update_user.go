package userapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ExecuteUpdateInput struct {
	Name     string    `validate:"gte=3"`
	Email    string    `validate:"email"`
	Password string    `validate:"gte=8"`
	ID       uuid.UUID `validate:"required"`
	ReqID    uuid.UUID `validate:"required"`
}

func (e ExecuteUpdateInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "Name":
					return fmt.Errorf("%w: %s", userdomain.ErrName, e.Error())
				case "Email":
					return fmt.Errorf("%w: %s", userdomain.ErrEmail, e.Error())
				case "Password":
					return fmt.Errorf("%w: %s", userdomain.ErrPassword, e.Error())
				case "ID":
					return fmt.Errorf("%w: %s", userdomain.ErrID, e.Error())
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

type ExecuteUpdateOutput struct {
	
}

func (u *Usecase) ExecuteUpdate(ctx context.Context, arg ExecuteUpdateInput) ()
