package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
)

type ExecuteGenerateTokenPairInput struct {
	UserID string
}

type ExecuteGenerateTokenPairOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (u *Usecase) ExecuteGenerateTokenPair(ctx context.Context, arg ExecuteGenerateTokenPairInput) (ExecuteGenerateTokenPairOutput, error) {
	scope := "jwtapp.ExecuteGenerateTokenPair"
	reqID := middleware.GetReqID(ctx)

	if err := u.store.IncrementTokenVersion(ctx, reqID, arg.UserID); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteGenerateTokenPairOutput{}, jwtdomain.ErrInternal
	}

	accessToken, err := u.executeGenerateAccessToken(ctx, arg.UserID)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteGenerateTokenPairOutput{}, jwtdomain.ErrInternal
	}

	refreshToken, err := u.executeGenerateRefreshToken(ctx,arg.UserID)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteGenerateTokenPairOutput{}, jwtdomain.ErrInternal
	}

	return ExecuteGenerateTokenPairOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.cfg.AccessTTL),
	}, nil
}
