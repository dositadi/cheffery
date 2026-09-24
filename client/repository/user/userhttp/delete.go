package userhttp

import (
	"errors"
	"net/http"

	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	scope := "userhttp.Delete"
	reqID := middleware.GetReqID(r.Context())

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.PrintInfo(reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		customerror.UnprocessableEntityResponse(w, r, userdomain.ErrID.Error())
		return
	}

	if err := h.port.ExecuteDelete(r.Context(), id); err != nil {
		h.logger.PrintInfo(reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		switch {
		case errors.Is(err, userdomain.ErrUnauthorized):
			customerror.UnauthorizedResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrTimeout):
			customerror.RequestTimeoutResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrNotFound):
			customerror.NotFoundResponse(w, r, err.Error())
			return
		default:
			customerror.InternalResponse(w, r, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
