package errors

import (
	"errors"
)

var NotFound = errors.New("not found")

type NotFoundError struct {
	ProblemDetails
}

func (e *NotFoundError) Error() string {
	return e.ProblemDetails.Error()
}

func (e *NotFoundError) Unwrap() error {
	return NotFound
}
