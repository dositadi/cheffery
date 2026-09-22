package authhttp

import (
	"net/http"

	"github.com/dositadi/cheffery/client/platform/middleware"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	logger     logger.Logger
	authClient auth.IssuerClient
	usecase    Usecase
}

func New(logger logger.Logger, authClient auth.IssuerClient, usecase Usecase) *Handler {
	return &Handler{
		logger:     logger,
		authClient: authClient,
		usecase:    usecase,
	}
}

func (h *Handler) Handler() http.Handler {
	r := chi.NewRouter()

	platformMW := middleware.NewPlatformMW(h.logger, h.authClient)

	r.Post("/login", h.Login)
	r.With(platformMW.ValidateAccess).Post("/logout", h.Logout)
	r.With(platformMW.ValidateRefresh).Post("/rotate-refresh", h.RotateRefresh)

	return r
}
