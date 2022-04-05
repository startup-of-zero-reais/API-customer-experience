package validation

import (
	"errors"
)

type (
	// UnauthorizedError is a error to be returned when the user is unauthorized
	Unauthorized struct {
		Err error
	}
)

// UnauthorizedError returns an error to be returned when the user is unauthorized
func UnauthorizedError(err string) *Unauthorized {
	return &Unauthorized{
		Err: errors.New(err),
	}
}

// Error returns the error message
func (e *Unauthorized) Error() string {
	return e.Err.Error()
}
