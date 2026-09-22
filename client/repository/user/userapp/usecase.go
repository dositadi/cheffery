package userapp

import (
	"github.com/dositadi/cheffery/protoc_gen/protoc/repository"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-playground/validator/v10"
)

type Usecase struct {
	logger     logger.Logger
	repoClient repository.RepositoryClient
	validate   *validator.Validate
}

func New(logger logger.Logger, repoClient repository.RepositoryClient,validate   *validator.Validate) *Usecase {
	return &Usecase{
		logger:     logger,
		repoClient: repoClient,
		validate: validate,
	}
}
