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
