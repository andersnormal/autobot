package nats

import (
	"fmt"
)

// NewError is return a new error with
func NewError(format string, a ...interface{}) error {
	return &queueError{fmt.Sprintf(format, a)}
}

type queueError struct {
	err string
}

func (e *queueError) Error() string {
	return fmt.Sprintf("%s", e.err)
}
