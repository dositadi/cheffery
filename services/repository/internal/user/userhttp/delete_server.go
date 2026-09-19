package userhttp

import (
	"context"
	"errors"
	"fmt"

	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/repository/internal/user/userapp"
	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteUser(ctx context.Context, req *repository.DeleteUserRequest) (*repository.DeleteUserResponse, error) {
	reqID := req.GetReqID()
	if reqID == "" {
		reqID = fmt.Sprintf("delete-user:%s", uuid.NewString())
	}
	scope := "userhttp.DeleteServer"

	userId, err := uuid.Parse(req.GetId())
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

	if err := s.executor.ExecuteDelete(ctx, userapp.ExecuteDeleteInput{
		ID:    userId,
		ReqID: req.GetReqID(),
	}); err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, userdomain.ErrID) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, userdomain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, userdomain.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &repository.DeleteUserResponse{}, nil
}
