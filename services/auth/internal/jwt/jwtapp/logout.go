package jwtapp

import (
	"context"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
)

type ExecuteLogoutInput struct {
	AccessToken string
}

func (u *Usecase) ExecuteLogout(ctx context.Context, arg ExecuteLogoutInput) (string, error) {
	scope := "jwtapp.ExecuteLogout"
	reqID := middleware.GetReqID(ctx)

	claim, err := u.ExecuteValidateAccessToken(ctx, arg.AccessToken)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		return "", err
	}

	if err := u.store.BlacklistAccessToken(ctx, reqID, claim.ID, claim.ExpiresAt.Time); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", jwtdomain.ErrInternal
	}

	if err := u.store.RevokeAllRefreshToken(ctx, reqID, claim.Subject); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", jwtdomain.ErrInternal
	}

	return claim.Subject, nil
}
