package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (u *Usecase) executeGenerateRefreshToken(ctx context.Context, reqId, userId string) (string, error) {
	tokenId := uuid.NewString()
	scope := "jwtapp.ExecuteGenerateAccessToken()"

	claim := jwtdomain.JWTRefreshClaim{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtdomain.Issuer.String(),
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"cheffery:user", "cheffery:app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.cfg.AccessTTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenId,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(u.cfg.AccessKey))
	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", jwtdomain.ErrInternal
	}

	// Store the token
	if err := u.store.StoreRefreshToken(ctx, reqId, tokenId, userId, time.Now().Add(u.cfg.RefreshTTL)); err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", jwtdomain.ErrInternal
	}

	return token, nil
}
