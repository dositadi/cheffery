package userdomain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID `validate:"required, eq=16"`
	name         string    `validate:"required, gte=3"`
	email        string    `validate:"required, email"`
	passwordHash []byte    `validate:"required, gte=8"`
	version      int32     `validate:"required, gte=1"`
	createdAt    time.Time `validate:"required"`
	updatedAt    time.Time
}

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

type CompareFunc func(hashedPassword, password []byte) error

func (u User) Compare(password []byte, fn CompareFunc) error { return fn(u.passwordHash, password) }
