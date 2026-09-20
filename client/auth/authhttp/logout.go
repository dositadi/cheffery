package authhttp

import (
	"errors"
	"net/http"

	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/globalutil/httphelper"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/sharedkernel"
	"github.com/go-chi/chi/middleware"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	scope := "authhttp.Logout"
	reqID := middleware.GetReqID(r.Context())
	accessToken := httphelper.GetAccessToken(r)

	if err := h.usecase.Logout(r.Context(), sharedkernel.JWT(accessToken)); err != nil {
		h.logger.PrintInfo(reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		if errors.Is(err, authdomain.ErrUnauthorized) {
			customerror.UnauthorizedResponse(w, r, err.Error())
			return
		}
		if errors.Is(err, authdomain.ErrTimeout) {
			customerror.RequestTimeoutResponse(w, r, err.Error())
			return
		}
		customerror.InternalResponse(w, r, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
