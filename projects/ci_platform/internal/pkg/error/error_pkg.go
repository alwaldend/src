package pkgError

import (
	"errors"
	"net/http"
)

// thin error wrapper
type Error interface {
	// interface for errors.Unwrap
	Unwrap() error
	StatusCode() int
	error
}

type customError struct {
	statusCode int
	err        error
}

func New(err error) error {
	return NewWithCode(err, http.StatusInternalServerError)
}

func NewWithCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &customError{code, err}
}

func (customError *customError) Error() string {
	return customError.err.Error()
}

func (customError *customError) Unwrap() error {
	return errors.Unwrap(customError.err)
}

func (customError *customError) StatusCode() int {
	return customError.statusCode
}
