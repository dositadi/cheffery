package jwtdomain

import "errors"

var (
	ErrTokenRevoked     = errors.New("jwtdomain: token revoked")
	ErrTokenExpired     = errors.New("jwtdomain: token expired")
	ErrInvalidToken     = errors.New("jwtdomain: invalid token")
	ErrInvalidTokenType = errors.New("jwtdomain: invalid token type")
)
