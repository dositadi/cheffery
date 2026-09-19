package jwtserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) RotateRefreshTokenServer(ctx context.Context, req *auth.RotateRefreshTokenRequest) (*auth.RotateRefreshTokenResponse, error) {
	scope := "jwtserver.RotateRefreshTokenServer"
	reqID := req.GetRequestID()
	if reqID == "" {
		reqID = fmt.Sprintf("rotate-refresh:%s", uuid.NewString())
	}
	refreshToken := req.GetRequestID()

	tokenPair, err := s.port.ExecuteRotateRefreshToken(ctx, jwtapp.ExecuteRotateRefreshTokenInput{
		ReqID:        reqID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		if errors.Is(err, jwtdomain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &gen.RotateRefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    timestamppb.New(tokenPair.ExpiresAt),
	}, nil
}
