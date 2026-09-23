package userdomain

import (
	"errors"
)

var (
	ErrID       = errors.New("userdomain: invalid id")
	ErrName     = errors.New("userdomain: name should contain at least 3 characters")
	ErrEmail    = errors.New("userdomain: invalid email")
	ErrPassword = errors.New("userdomain: password should contain at least 8 ")
	ErrVersion  = errors.New("userdomain: version should be greater than one")
	ErrDomain   = errors.New("userdomain: ")
	ErrInternal = errors.New("userdomain: something went wrong")
)

var (
	ErrBadRequest    = errors.New("userdomain: bad request")
	ErrUnauthorized  = errors.New("userdomain: unauthorized, invalid credentials")
	ErrEmailConflict = errors.New("userdomain: email conflict")
	ErrTimeout       = errors.New("userdomain: timeout")
	ErrNotFound      = errors.New("userdomain: not found")
)
