package userhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/globalutil/httphelper"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GetResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	scope := "userhttp.Get"
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

	response, err := h.port.ExecuteGetByID(r.Context(), id)
	if err != nil {
		h.logger.PrintInfo(reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})

		switch {
		case errors.Is(err, userdomain.ErrNotFound):
			customerror.NotFoundResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrTimeout):
			customerror.RequestTimeoutResponse(w, r, err.Error())
			return
		default:
			customerror.InternalResponse(w, r, err.Error())
			return
		}
	}

	httphelper.RenderJSON(w, r, http.StatusFound, GetResponse{
		Name:      response.GetName(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt(),
		UpdatedAt: response.GetUpdatedAt(),
	})
}
