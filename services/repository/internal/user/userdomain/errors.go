package userdomain

import "errors"

var (
	ErrName          = errors.New("userdomain: name field should be at least 3 characters long")
	ErrEmail         = errors.New("userdomain: invalid email")
	ErrPassword      = errors.New("userdomain: password field should be at least 8 characters long")
	ErrVersion       = errors.New("userdomain: version field required")
	ErrCreatedAt     = errors.New("userdomain: created at field required")
	ErrInternal      = errors.New("userdomain: something wrong happened")
	ErrID            = errors.New("userdomain: invalid id")
	ErrReqID         = errors.New("userdomain: request id required")
	ErrEmailConflict = errors.New("userdomain: this email exists already")
	ErrTimeout       = errors.New("userdomain: timeout")
	ErrNotFound      = errors.New("userdomain: user not found")
)
