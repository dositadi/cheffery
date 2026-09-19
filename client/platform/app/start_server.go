package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dositadi/cheffery/services/shared/customerror"
)

func (a *App) startServer() {
	scope := "app.startServer"
	for _, conn := range a.conns {
		defer func() {
			if err := conn.Close(); err != nil {
				a.logger.PrintError(err, "close-grpc-conn", customerror.InternalError{
					Inner:   err,
					Message: err.Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
			}
		}()
	}

	server := http.Server{
		Addr:           toAddr(a.cfg.Server.Host, a.cfg.Server.Port),
		Handler:        a.router,
		MaxHeaderBytes: int(a.cfg.Server.MaxHeaderBytes),
		ReadTimeout:    a.cfg.Server.ReadTimeout,
		WriteTimeout:   a.cfg.Server.WriteTimeout,
		IdleTimeout:    a.cfg.Server.IdleTimeout,
	}

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	chErr := make(chan error)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			chErr <- err
		}
	}()

	select {
	case e := <-chErr:
		a.logger.PrintError(e, "server-start", customerror.InternalError{
			Inner:   e,
			Message: e.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		return
	case <-chSignal:
		// fall through
	}

	if err := server.Close(); err != nil {
		a.logger.PrintError(err, "server-start", customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	}
}
