package jwtdomain

import "github.com/golang-jwt/jwt/v5"

type JWTAccessClaim struct {
	UserID       string
	UserName     string
	TokenVersion int64
	jwt.RegisteredClaims
}

type JWTRefreshClaim struct {
	UserID       string
	TokenVersion int64
	jwt.RegisteredClaims
}

type JWTRegisteredClaimField string

func (j JWTRegisteredClaimField) String() string {
	return string(j)
}

const (
	Issuer JWTRegisteredClaimField = "cheffery.app"
)
