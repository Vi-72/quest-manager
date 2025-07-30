package errs

import (
	"errors"
	"fmt"
)

var ErrVersionIsInvalid = errors.New("version is invalid")

type VersionIsInvalidError struct {
	ParamName string
	Cause     error
}

func NewVersionIsInvalidError(paramName string, cause error) *VersionIsInvalidError {
	return &VersionIsInvalidError{
		ParamName: paramName,
		Cause:     cause,
	}
}

func NewVersionIsInvalidErrorWithCause(paramName string) *VersionIsInvalidError {
	return &VersionIsInvalidError{
		ParamName: paramName,
	}
}

func (e *VersionIsInvalidError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", ErrVersionIsInvalid, e.ParamName, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrVersionIsInvalid, e.ParamName)
}

func (e *VersionIsInvalidError) Unwrap() error {
	return ErrVersionIsInvalid
}
