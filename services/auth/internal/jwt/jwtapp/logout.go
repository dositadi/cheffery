package jwtapp

import (
	"context"

	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
)

type ExecuteLogoutInput struct {
	ReqID       string
	AccessToken string
}

func (u *Usecase) ExecuteLogout(ctx context.Context, arg ExecuteLogoutInput) (string, error) {
	scope := "jwtapp.ExecuteLogout"

	claim, err := u.ExecuteValidateAccessToken(ctx, arg.ReqID, arg.AccessToken)
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		return "", err
	}

	if err := u.store.BlacklistAccessToken(ctx, arg.ReqID, claim.ID, claim.ExpiresAt.Time); err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return "", jwtdomain.ErrInternal
	}

	if err := u.store.RevokeAllRefreshToken(ctx, arg.ReqID, claim.Subject); err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
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
