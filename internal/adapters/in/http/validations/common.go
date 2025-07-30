package validations

import (
	"errors"
	"quest-manager/internal/adapters/in/http/problems"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrValidationFailed = errors.New("validation failed")
)

// ValidationError представляет ошибку валидации с множественными деталями
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

// ConvertValidationErrorToProblem конвертирует ValidationError в RFC 7807 Problem Details
func ConvertValidationErrorToProblem(err *ValidationError) *problems.BadRequest {
	return problems.NewBadRequest(err.Error())
}

// ValidateBody проверяет что body запроса не nil
func ValidateBody(body interface{}, bodyName string) *ValidationError {
	if body == nil {
		return NewValidationError(bodyName, "is required")
	}
	return nil
}

// ValidateNotEmpty проверяет что строка не пустая
func ValidateNotEmpty(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fieldName, "is required and cannot be empty")
	}
	return nil
}

// ValidateUUID проверяет что строка является валидным UUID
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

// TrimAndValidateString обрезает пробелы и проверяет что строка не пустая
func TrimAndValidateString(value, fieldName string) (string, *ValidationError) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", NewValidationError(fieldName, "is required and cannot be empty")
	}
	return trimmed, nil
}
