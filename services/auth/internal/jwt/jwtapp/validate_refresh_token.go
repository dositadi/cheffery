package jwtapp

import (
	"context"
	"errors"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func (u *Usecase) ExecuteValidateRefreshToken(ctx context.Context, refreshToken string) (*jwtdomain.JWTRefreshClaim, error) {
	reqId := middleware.GetReqID(ctx)
	scope := "jwtapp.ExecuteValidateAccessToken"

	token, err := jwt.ParseWithClaims(refreshToken, jwtdomain.JWTRefreshClaim{}, func(t *jwt.Token) (any, error) {
		return []byte(u.cfg.RefreshKey), nil
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
			return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrTokenExpired
		}
		return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrInvalidToken
	}

	claim, ok := token.Claims.(*jwtdomain.JWTRefreshClaim)
	if !ok || !token.Valid {
		return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrInvalidTokenType
	}

	tokenPayload, err := u.store.GetRefreshToken(ctx, reqId, claim.ID)
	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrInternal
	}

	if tokenPayload.Revoked {
		return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrTokenRevoked
	}

	if time.Now().Compare(tokenPayload.ExpiresAt) < 1 {
		return &jwtdomain.JWTRefreshClaim{}, jwtdomain.ErrTokenExpired
	}

	return claim, nil
}
