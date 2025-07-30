package errs

import (
	"fmt"
	"net/http"
)

// ErrorWithStatus wraps an error and provides an associated HTTP status code
type ErrorWithStatus struct {
	Err        error
	StatusCode int
	Message    string
}

func (e *ErrorWithStatus) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *ErrorWithStatus) Unwrap() error {
	return e.Err
}

func NewInternalServerError(msg string, err error) *ErrorWithStatus {
	return &ErrorWithStatus{
		Err:        fmt.Errorf(msg, err),
		StatusCode: http.StatusInternalServerError,
	}
}
