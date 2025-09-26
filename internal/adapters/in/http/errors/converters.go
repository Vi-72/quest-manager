package errors

import "quest-manager/internal/pkg/errs"

func NewDomainValidationProblem(err *errs.DomainValidationError) *BadRequest {
	detail := "validation failed: field '" + err.Field + "' " + err.Message
	if err.Cause != nil {
		detail += " (cause: " + err.Cause.Error() + ")"
	}
	return NewBadRequest(detail)
}

func NewNotFoundProblem(err *errs.NotFoundError) *ProblemDetails {
	detail := err.Resource + " with id '" + err.ID + "' not found"
	if err.Cause != nil {
		detail += " (cause: " + err.Cause.Error() + ")"
	}

	return &ProblemDetails{
		Type:   "not-found",
		Title:  "Not Found",
		Status: 404,
		Detail: detail,
	}
}
