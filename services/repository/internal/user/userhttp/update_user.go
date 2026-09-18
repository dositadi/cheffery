package userhttp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/repository/internal/user/userapp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) UpdateUser(ctx context.Context, req *repository.UpdateUserRequest) (*repository.UpdateUserResponse, error) {
	reqID := req.ReqID
	scope := "userhttp.UpdateUser"

	userId, err := uuid.Parse(req.Id)
	if err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, status.Error(codes.Unauthenticated, userdomain.ErrID.Error())
	}

	response, err := s.executor.ExecuteUpdate(ctx, userapp.ExecuteUpdateInput{
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.NewPassword,
		OldPassword: req.OldPassword,
		ID:          userId,
		ReqID:       reqID,
	})
	if err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		if errors.Is(err, userdomain.ErrName) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, userdomain.ErrEmail) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, userdomain.ErrPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, userdomain.ErrID) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, userdomain.ErrEmailConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, userdomain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, userdomain.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &repository.UpdateUserResponse{
		Datachanged: response.DataChanged,
		Fields:      response.Fields,
		Timestamp:   timestamppb.New(response.Timestamp),
	}, nil
}
