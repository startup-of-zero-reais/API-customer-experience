package validation

import "errors"

type (
	// BadRequest is a struct to handle bad request error
	BadRequest struct {
		Err error
	}
)

// BadRequestError is a function to handle bad request error
func BadRequestError(err string) *BadRequest {
	return &BadRequest{
		Err: errors.New(err),
	}
}

// Error is a function implements error interface
func (e *BadRequest) Error() string {
	return e.Err.Error()
}
