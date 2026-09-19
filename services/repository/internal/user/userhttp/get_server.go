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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetUser(ctx context.Context, req *repository.GetUserRequest) (*repository.GetUserResponse, error) {
	reqID := req.GetReqID()
	if reqID == "" {
		reqID = fmt.Sprintf("get-user:%s", uuid.NewString())
	}
	scope := "userhttp.GetUser"

	userID, err := uuid.Parse(req.GetUserID())
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

	user, err := s.executor.ExecuteGet(ctx, userapp.ExecuteGetInput{
		ReqID:  reqID,
		UserID: userID,
	})
	if err != nil {
		s.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		if errors.Is(err, userdomain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, userdomain.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &repository.GetUserResponse{
		Id:        user.GetID().String(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		Version:   user.GetVersion(),
		CreatedAt: timestamppb.New(user.GetCreatedAt()),
		UpdatedAt: timestamppb.New(user.GetUpdatedAt()),
	}, nil
}
