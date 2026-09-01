package jwtserver

import (
	"context"
	"errors"

	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ValidateRefreshTokenServer(ctx context.Context, req *gen.ValidateRefreshTokenRequest) (*gen.ValidateRefreshTokenResponse, error) {
	reqId := req.GetRequestID()
	token := req.GetRefreshToken()
	scope := "jwtserver.ValidateRefreshTokenServer"

	claim, err := s.port.ExecuteValidateRefreshToken(ctx, reqId, token)
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

	return &gen.ValidateRefreshTokenResponse{
		UserID: claim.Subject,
	}, nil
}
