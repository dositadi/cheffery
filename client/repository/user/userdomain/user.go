package userdomain

import (
	"errors"
	"time"
	"uuid"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
)

type User struct {
	id           uuid.UUID `validate:"required, uuid"`
	name         string    `validate:"required, gte=3"`
	email        string    `validate:"required, email"`
	passwordHash []byte    `validate:"required, len=8, alphanum"`
	version      int32     `validate:"required, gte=1"`
	createdAt    time.Time `validate:"required"`
	updatedAt    time.Time `validate:"omitempty"`
}

var validate *validator.Validate

func (u *User) validate() error {
	validate = validator.New(validator.WithRequiredStructEnabled(), validator.WithPrivateFieldValidation())

	if err := validate.Struct(u); err != nil {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "id":
					return customerror.WrapValidateErr(ErrID, e.Error())
				case "name":
					return customerror.WrapValidateErr(ErrName, e.Error())
				case "email":
					return customerror.WrapValidateErr(ErrEmail, e.Error())
				case "password":
					return customerror.WrapValidateErr(ErrPassword, e.Error())
				case "version":
					return customerror.WrapValidateErr(ErrVersion, e.Error())
				case "createdAt":
					return customerror.WrapValidateErr(ErrDomain, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(ErrInternal, err.Error())
	}
	return nil
}

func New(id uuid.UUID, name string, email string, passwordHash []byte, version int32, createdAt time.Time, updatedAt time.Time) (*User, error) {
	user := &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		version:      version,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u User) GetID() uuid.UUID        { return u.id }
func (u User) GetName() string         { return u.name }
func (u User) GetEmail() string        { return u.email }
func (u User) GetPasswordHash() []byte { return u.passwordHash }
func (u User) GetVersion() int32       { return u.version }
func (u User) GetCreatedAt() time.Time { return u.createdAt }
func (u User) GetUpdatedAt() time.Time { return u.updatedAt }
