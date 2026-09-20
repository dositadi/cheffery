package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
)

type ExecuteRotateRefreshTokenInput struct {
	RefreshToken string
}

type ExecuteRotateRefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (u *Usecase) ExecuteRotateRefreshToken(ctx context.Context, arg ExecuteRotateRefreshTokenInput) (ExecuteRotateRefreshTokenOutput, error) {
	scope := "jwtapp.ExecuteRotateRefreshToken"
	reqID := middleware.GetReqID(ctx)

	claim, err := u.ExecuteValidateRefreshToken(ctx, arg.RefreshToken)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		return ExecuteRotateRefreshTokenOutput{}, err
	}

	if err := u.store.RevokeRefreshToken(ctx, reqID, claim.ID); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		return ExecuteRotateRefreshTokenOutput{}, jwtdomain.ErrInternal
	}

	tokenPair, err := u.ExecuteGenerateTokenPair(ctx, ExecuteGenerateTokenPairInput{
		UserID: claim.Subject,
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteRotateRefreshTokenOutput{}, err
	}

	return ExecuteRotateRefreshTokenOutput(tokenPair), nil
}
