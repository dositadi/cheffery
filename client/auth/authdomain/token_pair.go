package authdomain

import (
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
)

type TokenPair struct {
	accessToken  string    `validate:"required, jwt"`
	refreshToken string    `validate:"required, jwt"`
	expiresAt    time.Time `validate:"required"`
}

func (t TokenPair) validate(validate *validator.Validate) error {
	if err := validate.Struct(t); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			for _, e := range vErrs {
				switch e.StructField() {
				case "accessToken":
					return customerror.WrapValidateErr(ErrAccessJWT, e.Error())
				case "refreshToken":
					return customerror.WrapValidateErr(ErrRefreshJWT, e.Error())
				case "expiresAt":
					return customerror.WrapValidateErr(ErrDomain, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(ErrInternal, err.Error())
	}
	return nil
}

func NewTokenPair(accessToken string, refreshToken string, expiresAt time.Time, validate *validator.Validate) (*TokenPair, error) {
	pair := &TokenPair{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
	}
	if err := pair.validate(validate); err != nil {
		return nil, err
	}
	return pair, nil
}
