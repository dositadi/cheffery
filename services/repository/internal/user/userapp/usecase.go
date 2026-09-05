package userapp

import (
	"github.com/dositadi/cheffery/services/repository/internal/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-playground/validator/v10"
)

type Usecase struct {
	logger   logger.Logger
	repo     userdomain.Repository
	validate *validator.Validate
	bcrypt   Bcrypt
}

func New(logger logger.Logger, repo userdomain.Repository, validate *validator.Validate, bcrypt Bcrypt) *Usecase {
	return &Usecase{
		logger:   logger,
		repo:     repo,
		validate: validate,
		bcrypt:   bcrypt,
	}
}
