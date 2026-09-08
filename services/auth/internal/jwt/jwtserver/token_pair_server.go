package jwtserver

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/auth/internal/jwt/jwtapp"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) TokenPairServer(ctx context.Context, req *auth.GenerateTokenPairRequest) (*auth.GenerateTokenPairResponse, error) {
	reqId := req.GetRequestID()
	userId := req.GetUserID()
	scope := "jwtserver.TokenPairServer"

	response, err := s.port.ExecuteGenerateTokenPair(ctx, jwtapp.ExecuteGenerateTokenPairInput{
		UserID: userId,
		ReqID:  reqId,
	})
	if err != nil {
		s.logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GenerateTokenPairResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    timestamppb.New(time.Now().Add(s.cfg.AccessTTL)),
	}, nil
}
