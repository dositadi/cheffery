package middleware

import (
	"fmt"
	"net/http"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/go-chi/chi/middleware"
)

func (m *PlatformMiddleware) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		scope := "middleware.Recover"

		defer func() {
			if rcv := recover(); rcv != nil {
				message, ok := rcv.(string)
				if !ok {
					message = "something wrong happened"
				}

				m.logger.PrintError(nil, reqID, customerror.InternalError{
					Inner:   nil,
					Message: fmt.Sprintln("Recoverer middleware: %s", message),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})

				customerror.InternalResponse(w, r, message)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
