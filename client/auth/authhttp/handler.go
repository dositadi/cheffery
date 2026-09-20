package authhttp

import (
	"net/http"

	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-chi/chi/v5"
)

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

func (h *Handler) Handler() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	r.Post("/rotate-refresh", h.RotateRefresh)

	return r
}
