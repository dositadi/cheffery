package userhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/dositadi/cheffery/client/repository/user/userapp"
	"github.com/dositadi/cheffery/client/repository/user/userdomain"
	"github.com/dositadi/cheffery/globalutil/httphelper"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type CreatePayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	scope := "userhttp.Create"

	payload := httphelper.RequestPayload[CreatePayload](r.Body, h.logger, reqID)

	response, err := h.port.ExecuteCreate(r.Context(), userapp.ExecuteCreateInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
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
		case errors.Is(err, userdomain.ErrEmailConflict):
			customerror.ConflictResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrTimeout):
			customerror.RequestTimeoutResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrInternal):
			customerror.InternalResponse(w, r, err.Error())
			return
		default:
			customerror.BadRequestResponse(w, r, err.Error())
			return
		}
	}

	httphelper.RenderJSON(w, r, http.StatusCreated, CreateResponse{
		ID:        response.ID,
		CreatedAt: response.CreatedAt,
	})
}
