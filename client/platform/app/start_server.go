package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dositadi/cheffery/services/shared/customerror"
)

// start server
func (a *App) startServer() {
	scope := "app.startServer"
	reqID := "server-start"

	defer func() {
		for _, conn := range a.conns {
			if err := conn.Close(); err != nil {
				a.logger.PrintError(err, "close-grpc-conn", customerror.InternalError{
					Inner:   err,
					Message: err.Error(),
					Misc:    nil,
				}.Error(), map[string]string{
					"Context": scope,
				})
			}
		}
	}()

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
		a.logger.PrintInfo(reqID, fmt.Sprintln("Server running at http://%s", toAddr(a.cfg.Server.Host, a.cfg.Server.Port)), map[string]string{
			"Context": scope,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			chErr <- err
		}
	}()

	select {
	case e := <-chErr:
		a.logger.PrintFatal(e, "server-start", customerror.InternalError{
			Inner:   e,
			Message: e.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
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
