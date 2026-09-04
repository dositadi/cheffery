package customerror

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type InternalError struct {
	Inner   error
	Message string
	Misc    map[string]string
}

func (i InternalError) Unwrap() error {
	return i.Inner
}

func (i InternalError) Wrap(err error) error {
	return fmt.Errorf("%w, %w", i.Inner, err)
}

func (i InternalError) Error() string {
	return fmt.Sprintf("%v: %s|%+v", i.Inner, i.Message, i.Misc)
}

func RetryableError(err error) bool {
	if nil == err {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.SerializationFailure, pgerrcode.DeadlockDetected, pgerrcode.CannotConnectNow:
			return true
		}
		if pgerrcode.IsConnectionException(pgErr); err != nil {
			return true
		}
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return false
}
