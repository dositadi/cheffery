package app

import (
	"fmt"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/services/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Establish a conn with a grpc server
func establishConn(logger logger.Logger, addr string) *grpc.ClientConn {
	reqID := "establish-connection"
	scope := "app.initAuthClient"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.PrintFatal(err, reqID, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
	}

	return conn
}

func toAddr(host string, port uint16) string {
	return fmt.Sprintf("%s:%v", host, port)
}
