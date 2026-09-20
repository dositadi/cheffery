package authhttp

import "github.com/dositadi/cheffery/services/shared/logger"

type Handler struct {
	logger  logger.Logger
	usecase Usecase
}

func New(logger logger.Logger, usecase Usecase) *Handler {
	return &Handler{
		logger:  logger,
		usecase: usecase,
	}
}
