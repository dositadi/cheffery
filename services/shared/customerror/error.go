package customerror

import (
	"fmt"
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
