package httphelper

import (
	"encoding/json/v2"
	"io"

	"github.com/dositadi/cheffery/services/shared/customerror"
	"github.com/dositadi/cheffery/services/shared/logger"
)

func RequestPayload[T any](body io.ReadCloser, logger logger.Logger, reqId string) T {
	scope := "httphelper.RequestPayload"

	defer func() {
		if err := body.Close(); err != nil {
			logger.PrintError(err, reqId, customerror.InternalError{
				Inner:   err,
				Message: err.Error(),
				Misc:    nil,
			}.Error(), map[string]string{
				"Context": scope,
			})
		}
	}()

	buffer, err := io.ReadAll(body)
	if err != nil {
		logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		panic("failed to read request body")
	}

	var out T

	if err := json.Unmarshal(buffer, &out); err != nil {
		logger.PrintError(err, reqId, customerror.InternalError{
			Inner:   err,
			Message: err.Error(),
			Misc:    nil,
		}.Error(), map[string]string{
			"Context": scope,
		})
		panic("failed to unmarshal request body")
	}

	return out
}
