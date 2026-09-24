package customerror

import (
	"net/http"

	"github.com/go-chi/render"
)

func InternalResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, message)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	render.JSON(w, r, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusNotFound)
	render.JSON(w, r, message)
}

func RequestTimeoutResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusRequestTimeout)
	render.JSON(w, r, message)
}

func GatewayTimeoutResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusGatewayTimeout)
	render.JSON(w, r, message)
}

func ConflictResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusConflict)
	render.JSON(w, r, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, message)
}
