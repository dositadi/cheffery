package jwtapp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/golang-jwt/jwt/v5"
)

func (u *Usecase) ExecuteValidateAccessToken(ctx context.Context, reqId, tokenString string) (*jwtdomain.JWTAccessClaim, error) {
	scope := "jwtapp.ExecuteValidateAccessToken"

	token, err := jwt.ParseWithClaims(tokenString, jwtdomain.JWTAccessClaim{}, func(t *jwt.Token) (any, error) {
		return []byte(u.cfg.AccessKey), nil
	},
		jwt.WithIssuedAt(),
		jwt.WithIssuer(string(jwtdomain.Issuer)),
		jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, jwt.ErrTokenExpired) {
			return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrTokenExpired
		}
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrInvalidToken
	}

	claim, ok := token.Claims.(*jwtdomain.JWTAccessClaim)
	if !ok || !token.Valid {
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrInvalidTokenType
	}

	currVersion, err := u.store.GetTokenVersion(ctx, reqId, claim.UserID)
	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrInvalidToken
	}

	if claim.TokenVersion != currVersion {
		if err = u.store.RevokeAllRefreshToken(ctx, reqId, claim.Subject); err != nil {
			u.logger.PrintError(err, reqId, customerror.InternalError{
				Inner:   err,
				Message: err.Error(),
				Misc:    nil,
			}.Error(), map[string]string{
				"Context": scope,
			})
		}
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrInvalidToken
	}

	blacklisted, err := u.store.AccessTokenIsBlackListed(ctx, reqId, claim.ID)
	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrInternal
	}

	if blacklisted {
		return &jwtdomain.JWTAccessClaim{}, jwtdomain.ErrTokenBlacklisted
	}

	return claim, nil
}
