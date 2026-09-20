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

func (s *Server) GetUserByEmail(ctx context.Context, req *repository.GetUserByEmailRequest) (*repository.GetUserByEmailResponse, error) {
	reqID := middleware.GetReqID(ctx)
	scope := "userhttp.GetUser"

	user, err := s.executor.ExecuteGetByEmail(ctx, userapp.ExecuteGetByEmailInput{
		Email: req.GetEmail(),
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

	return &repository.GetUserByEmailResponse{
		Id:           user.GetID().String(),
		Name:         user.GetName(),
		Email:        user.GetEmail(),
		PasswordHash: user.GetPasswordHash(),
		Version:      user.GetVersion(),
		CreatedAt:    timestamppb.New(user.GetCreatedAt()),
		UpdatedAt:    timestamppb.New(user.GetUpdatedAt()),
	}, nil
}
