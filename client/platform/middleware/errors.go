package middleware

import "errors"

var (
	ErrInternal     = errors.New("middleware: something wrong happened")
	ErrUnauthorized = errors.New("middleware: unauthorized, invalid credentials")
	ErrTimeout      = errors.New("middleware: timeout")
)
