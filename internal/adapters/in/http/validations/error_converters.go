package validations

import (
	"quest-manager/internal/adapters/in/http/problems"
	"quest-manager/internal/pkg/errs"
)

// ConvertDomainValidationErrorToProblem converts DomainValidationError to Problem Details (400 Bad Request)
func ConvertDomainValidationErrorToProblem(err *errs.DomainValidationError) *problems.BadRequest {
	detail := "validation failed: field '" + err.Field + "' " + err.Message
	if err.Cause != nil {
		detail += " (cause: " + err.Cause.Error() + ")"
	}
	return problems.NewBadRequest(detail)
}

// ConvertNotFoundErrorToProblem converts NotFoundError to Problem Details (404 Not Found)
func ConvertNotFoundErrorToProblem(err *errs.NotFoundError) *problems.ProblemDetails {
	detail := err.Resource + " with id '" + err.ID + "' not found"
	if err.Cause != nil {
		detail += " (cause: " + err.Cause.Error() + ")"
	}

	return &problems.ProblemDetails{
		Type:   "not-found",
		Title:  "Not Found",
		Status: 404,
		Detail: detail,
	}
}
