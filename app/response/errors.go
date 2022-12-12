package response

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrBadRequest          = errors.New("bad request, check param or body")
	ErrUnprocessableEntity = errors.New("can't process request, check param or body")
	ErrInternalServerError = errors.New("internal server error")
)

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrInternalServerError):
		return http.StatusInternalServerError
	default:
		wrapErr := &wrapErr{}
		if errors.As(err, wrapErr) {
			return wrapErr.StatusCode
		}

		return http.StatusInternalServerError
	}
}

// RespondError takes an `error` and a `customErr message` args
// to log the error to system and return to client
func RespondError(err error, customErr ...error) (int, Response) {
	if len(customErr) > 0 {
		return getStatusCode(err), Response{Message: customErr[0].Error()}
	}

	return getStatusCode(err), Response{Message: err.Error()}
}

func RespondValidationError(err error, errors map[string]interface{}) (int, Response) {
	return getStatusCode(err), Response{Message: err.Error(), Errors: errors}
}

type wrapErr struct {
	StatusCode int
	Err        error
}

// implements error interface
func (e wrapErr) Error() string {
	return e.Err.Error()
}

// Unwrap Implements the errors.Unwrap interface
func (e wrapErr) Unwrap() error {
	return e.Err // Returns inner error
}

func WrapError(err error, statusCode int) error {
	return wrapErr{
		Err:        err,
		StatusCode: statusCode,
	}
}
