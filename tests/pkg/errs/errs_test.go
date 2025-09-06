package errs_test

import (
	"errors"
	"net/http"
	"testing"

	"quest-manager/internal/pkg/errs"

	"github.com/stretchr/testify/assert"
)

func TestDomainValidationError(t *testing.T) {
	err := errs.NewDomainValidationError("field", "is invalid")
	assert.EqualError(t, err, "domain validation error: field 'field' is invalid")

	cause := errors.New("cause")
	errWithCause := errs.NewDomainValidationErrorWithCause("field", "is invalid", cause)
	assert.EqualError(t, errWithCause, "domain validation error: field 'field' is invalid (cause: cause)")
}

func TestNotFoundError(t *testing.T) {
	err := errs.NewNotFoundError("Resource", "123")
	assert.EqualError(t, err, "Resource with id '123' not found")

	cause := errors.New("cause")
	errWithCause := errs.NewNotFoundErrorWithCause("Resource", "123", cause)
	assert.EqualError(t, errWithCause, "Resource with id '123' not found (cause: cause)")
}

func TestErrorWithStatus(t *testing.T) {
	baseErr := errors.New("boom")
	e := &errs.ErrorWithStatus{Err: baseErr, StatusCode: http.StatusBadRequest}
	assert.Equal(t, "boom", e.Error())
	assert.Equal(t, baseErr, errors.Unwrap(e))

	custom := &errs.ErrorWithStatus{Err: baseErr, StatusCode: http.StatusBadRequest, Message: "oops"}
	assert.Equal(t, "oops", custom.Error())
}

func TestNewInternalServerError(t *testing.T) {
	e := errs.NewInternalServerError("%s happened", errors.New("boom"))
	assert.Equal(t, http.StatusInternalServerError, e.StatusCode)
}

func TestValueIsRequiredError(t *testing.T) {
	err := errs.NewValueIsRequiredError("param")
	assert.EqualError(t, err, "value is required: param")
	assert.Equal(t, errs.ErrValueIsRequired, errors.Unwrap(err))

	cause := errors.New("cause")
	errWithCause := errs.NewValueIsRequiredErrorWithCause("param", cause)
	assert.EqualError(t, errWithCause, "value is required: param (cause: cause)")
	assert.Equal(t, errs.ErrValueIsRequired, errors.Unwrap(errWithCause))
}

func TestWrapInfrastructureError(t *testing.T) {
	cause := errors.New("db error")
	err := errs.WrapInfrastructureError("save failed", cause)
	assert.EqualError(t, err, "save failed: db error")
	assert.Equal(t, cause, errors.Unwrap(err))

	withoutCause := errs.WrapInfrastructureError("save failed", nil)
	assert.EqualError(t, withoutCause, "save failed")
}
