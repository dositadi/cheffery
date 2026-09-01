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

func (s *Server) ValidateAccessTokenServer(ctx context.Context, req *gen.ValidateAccessTokenRequest) (*gen.ValidateAccessTokenResponse, error) {
	reqId := req.GetRequestID()
	token := req.GetAccessToken()
	scope := "jwtserver.ValidateAccessTokenServer"

	claim, err := s.port.ExecuteValidateAccessToken(ctx, reqId,token)
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

	return &gen.ValidateAccessTokenResponse{
		UserID:       claim.Subject,
		TokenVersion: claim.TokenVersion,
	}, nil
}
