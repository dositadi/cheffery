package authapp

import (
	"context"
	"errors"

	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginInput struct {
	Email    string `validate:"required, email"`
	Password string `validate:"required, gte=8, alphanum"`
}

func (l LoginInput) validate(v *validator.Validate) error {
	if err := v.Struct(l); err != nil {
		if vErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			for _, e := range vErrs {
				switch e.StructField() {
				case "Email":
					return customerror.WrapValidateErr(authdomain.ErrEmail, e.Error())
				case "Password":
					return customerror.WrapValidateErr(authdomain.ErrPassword, e.Error())
				}
			}
		}
		return customerror.WrapValidateErr(authdomain.ErrInternal, err.Error())
	}
	return nil
}

func (u *Usecase) Login(ctx context.Context, arg LoginInput) (*authdomain.TokenPair, error) {
	reqID := middleware.GetReqID(ctx)
	scope := "authapp.Login"

	response, err := u.repoService.GetUserByEmail(ctx, &repository.GetUserByEmailRequest{
		Email: arg.Email,
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		st, ok := status.FromError(err)
		if !ok {
			return nil, authdomain.ErrInternal
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, authdomain.ErrUnauthorized
		case codes.DeadlineExceeded:
			return nil, authdomain.ErrTimeout
		case codes.Internal:
			return nil, authdomain.ErrInternal
		}
		return nil, authdomain.ErrInternal
	}

	if err := u.bcrypt.Compare(response.GetPasswordHash(), []byte(arg.Password)); err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, authdomain.ErrUnauthorized
	}

	tokenPair, err := u.authService.GenerateTokenPair(ctx, &auth.GenerateTokenPairRequest{
		UserID: response.GetId(),
	})
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		st, ok := status.FromError(err)
		if !ok {
			return nil, authdomain.ErrInternal
		}

		switch st.Code() {
		case codes.DeadlineExceeded:
			return nil, authdomain.ErrTimeout
		case codes.Internal:
			return nil, authdomain.ErrInternal
		}
		return nil, authdomain.ErrInternal
	}

	out, err := authdomain.NewTokenPair(tokenPair.AccessToken, tokenPair.RefreshToken, tokenPair.ExpiresAt.AsTime(), u.validate)
	if err != nil {
		u.logger.PrintError(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return nil, authdomain.ErrInternal
	}


	return out, nil
}
