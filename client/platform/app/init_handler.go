package app

import (
	"fmt"
	"time"

	"github.com/dositadi/cheffery/client/auth/authapp"
	"github.com/dositadi/cheffery/client/auth/authhttp"
	m "github.com/dositadi/cheffery/client/platform/middleware"
	"github.com/dositadi/cheffery/client/repository/user/userapp"
	"github.com/dositadi/cheffery/client/repository/user/userhttp"
	"github.com/dositadi/cheffery/globalutil/bcrypt/bcryptapp"
	"github.com/go-chi/chi/middleware"
)

type version uint

const (
	v1 = 1
)

var (
	baseUrl = fmt.Sprintf("/v%v/cheffery", v1)
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

	a.mountAuth()
	a.mountUser()
}

func (a *App) mountAuth() {
	bcrypt := bcryptapp.New(a.logger)
	useCase := authapp.New(a.authClient, a.repoClient, a.validate, bcrypt, a.logger)

	handler := authhttp.New(a.logger, a.authClient, useCase)

	a.router.Mount(baseUrl+"/auth", handler.Handler())
}

func (a *App) mountUser() {
	usecase := userapp.New(a.logger, a.repoClient, a.validate)

	handler := userhttp.New(a.logger, usecase, a.authClient)

	a.router.Mount(baseUrl+"/user", handler.Handler())
}
