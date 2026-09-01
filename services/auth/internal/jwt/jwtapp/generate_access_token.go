package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ExecuteGenerateAccessTokenInput struct {
	UserID string
	Name   string
}

func (u *Usecase) ExecuteGenerateAccessToken(ctx context.Context, reqId string, arg ExecuteGenerateAccessTokenInput) (string, error) {
	tokenId := uuid.NewString()
	scope := "jwtapp.ExecuteGenerateAccessToken()"

	version, err := u.store.GetTokenVersion(ctx, reqId, arg.UserID)
	if err != nil {
		u.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", err
	}

	claim := jwtdomain.JWTAccessClaim{
		UserID:       arg.UserID,
		UserName:     arg.Name,
		TokenVersion: version,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtdomain.Issuer.String(),
			Subject:   arg.UserID,
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
		return "", err
	}

	return token, nil
}
