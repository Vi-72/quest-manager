package errs

import (
	"errors"
	"fmt"
)

var ErrOperationFailed = errors.New("operation failed")

type OperationFailedError struct {
	Reason string
	Cause  error
}

func NewOperationFailedError(reason string) *OperationFailedError {
	return &OperationFailedError{
		Reason: reason,
	}
}

func NewOperationFailedErrorWithCause(reason string, cause error) *OperationFailedError {
	return &OperationFailedError{
		Reason: reason,
		Cause:  cause,
	}
}

func (e *OperationFailedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", ErrOperationFailed, e.Reason, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrOperationFailed, e.Reason)
}

func (e *OperationFailedError) Unwrap() error {
	return ErrOperationFailed
}
