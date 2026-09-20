package authhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/dositadi/cheffery/client/auth/authapp"
	"github.com/dositadi/cheffery/client/auth/authdomain"
	"github.com/dositadi/cheffery/globalutil/httphelper"
	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	scope := "authhttp.Login"
	reqID := middleware.GetReqID(r.Context())

	payload := httphelper.RequestPayload[LoginPayload](r.Body, h.logger, reqID)

	resp, err := h.usecase.Login(r.Context(), authapp.LoginInput{
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
		if errors.Is(err, authdomain.ErrTimeout) {
			customerror.RequestTimeoutResponse(w, r, err.Error())
			return
		}
		customerror.InternalResponse(w, r, err.Error())
		return
	}

	httphelper.RenderJSON(w, r, http.StatusOK, LoginResponse{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresAt:    resp.GetExpiresAt(),
	})
}
