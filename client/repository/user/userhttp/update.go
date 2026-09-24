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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UpdatePayload struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
}

type UpdateResponse struct {
	DataChanged bool      `json:"dataChanged"`
	Fields      []string  `json:"fields,omitempty"`
	Timestamp   time.Time `json:"timeStamp"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	scope := "userhttp.Update"
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

	payload := httphelper.RequestPayload[UpdatePayload](r.Body, h.logger, reqID)

	response, err := h.port.ExecuteUpdate(r.Context(), userapp.ExecuteUpdateInput{
		Name:        payload.Name,
		Email:       payload.Email,
		Password:    payload.Password,
		OldPassword: payload.OldPassword,
		ID:          id,
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
		case errors.Is(err, userdomain.ErrNotFound):
			customerror.NotFoundResponse(w, r, err.Error())
			return
		case errors.Is(err, userdomain.ErrInternal):
			customerror.InternalResponse(w, r, err.Error())
			return
		default:
			customerror.BadRequestResponse(w, r, err.Error())
			return
		}
	}

	httphelper.RenderJSON(w, r, http.StatusOK, UpdateResponse{
		DataChanged: response.DataChanged,
		Fields:      response.Fields,
		Timestamp:   response.Timestamp,
	})
}
