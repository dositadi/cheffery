package jwtserver

import (
	"context"
	"errors"

	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) LogoutHandler(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error) {
	reqId := req.GetRequestID()
	accessToken := req.GetAccessToken()
	scope := "jwtserver.LogoutHandler"

	userID, err := s.port.ExecuteLogout(ctx, jwtapp.ExecuteLogoutInput{
		ReqID:       reqId,
		AccessToken: accessToken,
	})
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

	return &gen.LogoutResponse{
		UserID: userID,
	}, nil
}
