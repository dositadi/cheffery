package jwtapp

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExecuteGenerateTokenPairInput struct {
	UserID string
	ReqID  string
}

type ExecuteGenerateTokenPairOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (u *Usecase) ExecuteGenerateTokenPair(ctx context.Context, arg ExecuteGenerateTokenPairInput) (ExecuteGenerateTokenPairOutput, error) {
	scope := "jwtapp.ExecuteGenerateTokenPair"

	accessToken, err := u.executeGenerateAccessToken(ctx, arg.ReqID, arg.UserID)
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteGenerateTokenPairOutput{}, status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := u.executeGenerateRefreshToken(ctx, arg.ReqID, arg.UserID)
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteGenerateTokenPairOutput{}, status.Error(codes.Internal, err.Error())
	}

	return ExecuteGenerateTokenPairOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(u.cfg.AccessTTL),
	}, nil
}
