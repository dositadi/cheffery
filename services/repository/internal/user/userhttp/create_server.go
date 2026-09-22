package userhttp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/repository/internal/user/userapp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateUser(ctx context.Context, req *repository.CreateUserRequest) (*repository.CreateUserResponse, error) {
	reqID := middleware.GetReqID(ctx)
	scope := "userhttp.CreateUser"

	response, err := s.executor.ExecuteCreate(ctx, userapp.ExecuteCreateInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, userdomain.ErrName) || errors.Is(err, userdomain.ErrEmail) || errors.Is(err, userdomain.ErrPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, userdomain.ErrID) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, userdomain.ErrEmailConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, userdomain.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &repository.CreateUserResponse{
		Id:        response.ID.String(),
		CreatedAt: timestamppb.New(response.CreatedAt),
	}, nil
}
