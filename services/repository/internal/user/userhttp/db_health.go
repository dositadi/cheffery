package userhttp

import (
	"context"
	"fmt"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) DBHealth(ctx context.Context, req *repository.DBHealthRequest) (*repository.DBHealthResponse, error) {
	reqID := req.ReqID
	if reqID == "" {
		reqID = fmt.Sprintf("db-health:%s", uuid.NewString())
	}
	scope := "userhttp.DBHealth"

	response := s.healthChecker.Check(ctx)
	if nil == response {
		s.logger.PrintError(nil, reqID, customerror.InternalError{
			Inner:   nil,
			Message: "Db check is nil",
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, status.Error(codes.Unavailable, "Db check is unavailable")
	}

	return &repository.DBHealthResponse{
		Status: string(response.Status),
		DbHealth: &repository.DBHealth{
			Connected:      response.DBHealth.Connected,
			PoolHealthy:    response.DBHealth.PoolHealthy,
			TotalConns:     response.DBHealth.TotalConns,
			IdleConns:      response.DBHealth.IdleConns,
			AcquiredConns:  response.DBHealth.AcquiredConns,
			AcquiredCount:  response.DBHealth.AcquiredCount,
			CreatedCount:   response.DBHealth.CreatedCount,
			ReleasedCount:  response.DBHealth.ReleasedCount,
			DestroyedCount: response.DBHealth.DestroyedCount,
			Message:        response.DBHealth.Message,
		},
		Latency:   durationpb.New(response.Latency),
		Timestamp: timestamppb.New(response.Timestamp),
	}, nil
}
