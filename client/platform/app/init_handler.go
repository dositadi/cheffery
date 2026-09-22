package app

import (
	"fmt"
	"time"

	m "github.com/dositadi/cheffery/client/platform/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type version uint

const (
	v1 = 1
)

func (a *App) mountHandlers() {
	mw := m.NewPlatformMW(a.logger, a.authClient)

	// Set up global middlewares.
	a.router.Use(mw.Recoverer, middleware.CleanPath, middleware.Timeout(a.cfg.Server.WriteTimeout), middleware.RequestID, middleware.Logger, middleware.ThrottleWithOpts(middleware.ThrottleOpts{
		Limit:        500,
		BacklogLimit: 300,
		RetryAfterFn: func(ctxDone bool) time.Duration {
			if !ctxDone {
				// Note: this is the time when the client should retry again
				return a.cfg.Server.WriteTimeout
			}
			return time.Duration(0)
		},
		BacklogTimeout: a.cfg.Server.WriteTimeout,
	}))

	a.router.Mount(fmt.Sprintf("/cheffery/v%v", v1), a.router.Group(func(r chi.Router) {
		// Mount all domain handlers here
	}))
}
