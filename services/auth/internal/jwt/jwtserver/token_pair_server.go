package jwtserver

import (
	"context"

	gen "github.com/dositadi/cheffery/protoc_gen/protoc/auth"
)

func (s *Server) TokenPairServer(ctx context.Context, req *gen.GenerateTokenPairRequest) (*gen.GenerateTokenPairResponse, error) {
	
}