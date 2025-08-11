package validations

import (
	"errors"
	"strings"

	"quest-manager/internal/adapters/in/http/problems"

	"github.com/google/uuid"
)

var (
	ErrValidationFailed = errors.New("validation failed")
)

// ValidationError represents a validation error with multiple details
type ValidationError struct {
	Message string
	Field   string
	Cause   error
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewValidationErrorWithCause(field, message string, cause error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return ErrValidationFailed.Error() + ": field '" + e.Field + "' " + e.Message + " (cause: " + e.Cause.Error() + ")"
	}
	return ErrValidationFailed.Error() + ": field '" + e.Field + "' " + e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}

// ConvertValidationErrorToProblem converts ValidationError to RFC 7807 Problem Details
func ConvertValidationErrorToProblem(err *ValidationError) *problems.BadRequest {
	return problems.NewBadRequest(err.Error())
}

// ValidateBody checks that request body is not nil
func ValidateBody(body interface{}, bodyName string) *ValidationError {
	if body == nil {
		return NewValidationError(bodyName, "is required")
	}
	return nil
}

// ValidateNotEmpty checks that string is not empty
func ValidateNotEmpty(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fieldName, "is required and cannot be empty")
	}
	return nil
}

// ValidateUUID checks that string is a valid UUID
func ValidateUUID(value, fieldName string) (uuid.UUID, *ValidationError) {
	if value == "" {
		return uuid.UUID{}, NewValidationError(fieldName, "is required")
	}

	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, NewValidationErrorWithCause(fieldName, "must be a valid UUID format", err)
	}

	return parsedUUID, nil
}

// TrimAndValidateString trims spaces and checks that string is not empty
func TrimAndValidateString(value, fieldName string) (string, *ValidationError) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", NewValidationError(fieldName, "is required and cannot be empty")
	}
	return trimmed, nil
}
