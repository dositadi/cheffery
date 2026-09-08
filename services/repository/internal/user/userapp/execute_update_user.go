package userapp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/repository/internal/user/userpostgres"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ExecuteUpdateInput struct {
	Name        string    `validate:"gte=3"`
	Email       string    `validate:"email"`
	Password    string    `validate:"gte=8"`
	OldPassword string    `validate:"gte=8"`
	ID          uuid.UUID `validate:"required, eq=16"`
	ReqID       string    `validate:"required"`
}

func (e ExecuteUpdateInput) validate(validate *validator.Validate) error {
	if err := validate.Struct(e); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch e.StructField() {
				case "Name":
					return fmt.Errorf("%w: %s", userdomain.ErrName, e.Error())
				case "Email":
					return fmt.Errorf("%w: %s", userdomain.ErrEmail, e.Error())
				case "Password":
					return fmt.Errorf("%w: %s", userdomain.ErrPassword, e.Error())
				case "ID":
					return fmt.Errorf("%w: %s", userdomain.ErrID, e.Error())
				}
			}
		}
		return userdomain.ErrInternal
	}
	return nil
}

type ExecuteUpdateOutput struct {
	DataChanged bool
	Fields      []string
	UpdatedAt   time.Time
}

func (u *Usecase) ExecuteUpdate(ctx context.Context, arg ExecuteUpdateInput) (ExecuteUpdateOutput, error) {
	scope := "userapp.ExecuteUpdate"

	if err := arg.validate(u.validate); err != nil {
		if arg.ReqID == "" {
			arg.ReqID = "userapp.ExecuteCreate-Request"
		}
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteUpdateOutput{}, err
	}

	user, err := u.ExecuteGet(ctx, ExecuteGetInput{
		ReqID:  arg.ReqID,
		UserID: arg.ID,
	})
	if err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return ExecuteUpdateOutput{}, err
	}

	var changedFields []string

	if arg.Name != "" && arg.Name != user.GetName() {
		user.SetName(arg.Name)
		changedFields = append(changedFields, "name")
	}

	if arg.Email != "" && arg.Email != user.GetEmail() {
		user.SetEmail(arg.Email)
		changedFields = append(changedFields, "email")
	}

	if err := user.Compare([]byte(arg.OldPassword), u.bcrypt.Compare); err == nil {
		if err := user.Compare([]byte(arg.Password), u.bcrypt.Compare); err != nil {
			newHashedPassword, err := u.bcrypt.GenerateHash([]byte(arg.Password))
			if err != nil {
				u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
					Inner:   err,
					Message: err.Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
				return ExecuteUpdateOutput{}, userdomain.ErrInternal
			}
			user.SetPasswordHash(newHashedPassword)
			changedFields = append(changedFields, "password")
		}
	}

	if err := u.repo.UpdateUser(ctx, arg.ReqID, userpostgres.UpdateUserInput{
		Name:         user.GetName(),
		Email:        user.GetEmail(),
		PasswordHash: user.GetPasswordHash(),
		ID:           user.GetID(),
		Version:      user.GetVersion(),
	}); err != nil {
		u.logger.PrintError(err, arg.ReqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, userpostgres.ErrEmailConflict) {
			return ExecuteUpdateOutput{}, userdomain.ErrEmailConflict
		}

		if errors.Is(err, userpostgres.ErrRequestTimeout) {
			return ExecuteUpdateOutput{}, userdomain.ErrTimeout
		}
		return ExecuteUpdateOutput{}, userdomain.ErrInternal
	}

	return ExecuteUpdateOutput{
		DataChanged: len(changedFields) > 0,
		Fields:      changedFields,
		UpdatedAt:   time.Now(),
	}, nil
}
