package authapp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
)

type LoginInput struct {
	Email    string `validate:"required, email"`
	Password string `validate:"required, gte=8, alphanum"`
}

func (l LoginInput) validate(v *validator.Validate) error {
	if err := v.Struct(l); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			for _, e := range vErrs {
				switch e.StructField() {
				case "Email":
					return customerror.WrapValidateErr(authdomain.ErrEmail, e.Error())
				case "Password":
					return customerror.WrapValidateErr(authdomain.ErrPassword, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(authdomain.ErrInternal, err.Error())
	}
	return nil
}

func (u *Usecase) Login(ctx context.Context, arg LoginInput) (*authdomain.TokenPair, error) {
	reqID := middleware.GetReqID(ctx)
	scope := "authapp.Login"

	u.repoService.GetUser()
}
