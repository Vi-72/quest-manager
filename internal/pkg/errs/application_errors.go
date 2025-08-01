package errs

import "fmt"

// DomainValidationError represents validation error at domain/application level
type DomainValidationError struct {
	Field   string
	Message string
	Cause   error
}

func (e *DomainValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("domain validation error: field '%s' %s (cause: %v)", e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("domain validation error: field '%s' %s", e.Field, e.Message)
}

func NewDomainValidationError(field, message string) *DomainValidationError {
	return &DomainValidationError{
		Field:   field,
		Message: message,
	}
}

func NewDomainValidationErrorWithCause(field, message string, cause error) *DomainValidationError {
	return &DomainValidationError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// NotFoundError represents "not found" error
type NotFoundError struct {
	Resource string
	ID       string
	Cause    error
}

func (e *NotFoundError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s with id '%s' not found (cause: %v)", e.Resource, e.ID, e.Cause)
	}
	return fmt.Sprintf("%s with id '%s' not found", e.Resource, e.ID)
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

func NewNotFoundErrorWithCause(resource, id string, cause error) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Cause:    cause,
	}
}
