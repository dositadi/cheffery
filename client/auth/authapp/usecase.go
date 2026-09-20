package authapp

import (
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-playground/validator/v10"
)

type Usecase struct {
	logger      logger.Logger
	validate    *validator.Validate
	authService auth.IssuerClient
	repoService repository.RepositoryClient
	bcrypt      Bcrypt
}

func New(authService auth.IssuerClient, repoService repository.RepositoryClient, validate *validator.Validate, bcrypt Bcrypt, logger logger.Logger) *Usecase {
	return &Usecase{
		logger:      logger,
		validate:    validate,
		authService: authService,
		repoService: repoService,
		bcrypt:      bcrypt,
	}
}
