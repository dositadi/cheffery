package userdomain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID `validate:"required, eq=16"`
	name      string    `validate:"required, gte=3"`
	email     string    `validate:"required, email"`
	password  []byte    `validate:"required, gte=8"`
	version   int32     `validate:"required, gte=1"`
	createdAt time.Time `validate:"required"`
	updatedAt time.Time
}

var (
	ErrName      = errors.New("name field should be at least 3 characters long")
	ErrEmail     = errors.New("invalid email")
	ErrPassword  = errors.New("password field should be at least 8 characters long")
	ErrVersion   = errors.New("version field required")
	ErrCreatedAt = errors.New("created at field required")
	ErrInternal  = errors.New("Something wrong happened")
	ErrID        = errors.New("invalid id")
)

var validate *validator.Validate

func (u *User) validate() error {
	validate = validator.New(validator.WithRequiredStructEnabled(), validator.WithPrivateFieldValidation())

	if err := validate.Struct(u); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "id":
					return ErrID
				case "name":
					return ErrName
				case "email":
					return ErrEmail
				case "password":
					return ErrPassword
				case "version":
					return ErrVersion
				case "createdAt":
					return ErrCreatedAt
				}
			}
		}
		return ErrInternal
	}
	return nil
}

func New(id uuid.UUID, name string, email string, password []byte, version int32, createdAt time.Time, updatedAt time.Time) (*User, error) {
	user := &User{
		id:        id,
		name:      name,
		email:     email,
		password:  password,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}
