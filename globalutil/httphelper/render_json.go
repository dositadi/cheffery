package httphelper

import (
	"net/http"

	"github.com/go-chi/render"
)

func RenderJSON(w http.ResponseWriter, r *http.Request, statusCode int, message any) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, message)
}
