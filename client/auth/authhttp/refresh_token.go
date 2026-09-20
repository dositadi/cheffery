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

func (h *Handler) RotateRefresh(w http.ResponseWriter, r *http.Request) {
	scope := "authhttp.RotateRefresh"
	reqID := middleware.GetReqID(r.Context())
	refreshToken := httphelper.GetRefreshToken(r)

	resp, err := h.usecase.RotateRefresh(r.Context(), sharedkernel.JWT(refreshToken))
	if err != nil {
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

	httphelper.RenderJSON(w, r, http.StatusOK, TokenPairResponse{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresAt:    resp.GetExpiresAt(),
	})
}
