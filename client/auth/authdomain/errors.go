package authdomain

import "errors"

var (
	ErrAccessJWT  = errors.New("authdomain: invalid access token")
	ErrRefreshJWT = errors.New("authdomain: invalid refresh token")
	ErrDomain     = errors.New("authdomain: ")
)

var (
	ErrID           = errors.New("authdomain: invalid id")
	ErrTokenVersion = errors.New("authdomain: invalid version")
	ErrInternal     = errors.New("authdomain: something went wrong")
)

var (
	ErrEmail = errors.New("authdomain: invalid email")
	ErrPassword = errors.New("authdomain: email should be at least 8 chars long")
)
