package errs

import (
	"errors"
	"fmt"
)

// InfrastructureError is returned when low-level (DB, IO, network, etc.) operations fail.
type InfrastructureError struct {
	msg string
	err error
}

func (e *InfrastructureError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *InfrastructureError) Unwrap() error {
	return e.err
}

// WrapInfrastructureError creates a wrapped InfrastructureError
func WrapInfrastructureError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return &InfrastructureError{msg: msg, err: err}
}
