package userhttp

import (
	"net/http"

	"github.com/dositadi/cheffery/client/platform/middleware"
	"github.com/dositadi/cheffery/protoc_gen/protoc/auth"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	logger     logger.Logger
	port       UsecasePort
	authClient auth.IssuerClient
}

func New(logger logger.Logger, port UsecasePort, authClient auth.IssuerClient) *Handler {
	return &Handler{
		port:       port,
		logger:     logger,
		authClient: authClient,
	}
}

func (h *Handler) Handler() http.Handler {
	r := chi.NewRouter()

	platformMW := middleware.NewPlatformMW(h.logger, h.authClient)

	r.Post("/", h.Create)
	r.With(platformMW.ValidateAccess).Post("/{id}", h.Delete)
	r.With(platformMW.ValidateAccess).Patch("/{id}", h.Update)
	r.With(platformMW.ValidateAccess).Get("/{id}", h.Get)

	return r
}
