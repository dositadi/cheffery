package authdomain

import (
	"errors"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AccessCredentials struct {
	id           uuid.UUID `validate:"required, uuid"`
	tokenVersion int64     `validate:"required"`
}

func (a AccessCredentials) validate(validate *validator.Validate) error {
	if err := validate.Struct(a); err != nil {
		if vErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range vErrs {
				switch e.StructField() {
				case "id":
					return customerror.WrapValidateErr(ErrID, e.Error())
				case "tokenVersion":
					return customerror.WrapValidateErr(ErrTokenVersion, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(ErrInternal, err.Error())
	}
	return nil
}

func NewAccessCredentials(id uuid.UUID, tokenVersion int64, validate *validator.Validate) (*AccessCredentials, error) {
	credential := &AccessCredentials{
		id:           id,
		tokenVersion: tokenVersion,
	}

	if err := credential.validate(validate); err != nil {
		return nil, err
	}
	return credential, nil
}
