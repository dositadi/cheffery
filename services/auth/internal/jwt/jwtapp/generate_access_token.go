package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ExecuteGenerateAccessTokenInput struct {
	UserID string
	Name   string
}

func (u *Usecase) ExecuteGenerateAccessToken(ctx context.Context, arg ExecuteGenerateAccessTokenInput) string {
	tokenId := uuid.NewString()

	claim := jwtdomain.JWTClaim{
		UserID:   arg.UserID,
		UserName: arg.Name,
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

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(u.cfg.AccessKey)
	if err != nil {
		
	}

}
