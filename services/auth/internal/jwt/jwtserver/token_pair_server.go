package jwtserver

import (
	"context"
	"errors"
	"time"

	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) TokenPairServer(ctx context.Context, req *gen.GenerateTokenPairRequest) (*gen.GenerateTokenPairResponse, error) {
	reqId := req.GetRequestID()
	userId := req.GetUserID()
	scope := "jwtserver.TokenPairServer"

	accessToken, err := s.port.ExecuteGenerateAccessToken(ctx, reqId, userId)
	if err != nil {
		s.logger.PrintError(err, reqId, customerror.InternalError{
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

	refreshToken, err := s.port.ExecuteGenerateRefreshToken(ctx, reqId, userId)
	if err != nil {
		s.logger.PrintError(err, reqId, customerror.InternalError{
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

	return &gen.GenerateTokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(time.Now().Add(s.cfg.AccessTTL)),
	}, nil
}
