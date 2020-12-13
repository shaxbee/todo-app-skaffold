package httprouter

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Status  int
	Message string
	Cause   error
}

func NewError(status int, opts ...ErrorOpt) Error {
	err := Error{Status: status}

	for _, opt := range opts {
		opt(&err)
	}

	if err.Message == "" {
		err.Message = http.StatusText(status)
	}

	return err
}

func (e Error) Unwrap() error {
	return e.Cause
}

func (e Error) Error() string {
	return e.Message
}

type ErrorOpt func(*Error)

func Message(message string) ErrorOpt {
	return func(e *Error) {
		e.Message = message
	}
}

func Messagef(format string, a ...interface{}) ErrorOpt {
	return func(e *Error) {
		e.Message = fmt.Sprintf(format, a...)
	}
}

func Cause(cause error) ErrorOpt {
	return func(e *Error) {
		e.Cause = cause
	}
}

func IsStatus(err error, status int) bool {
	var httpErr Error
	switch {
	case errors.Is(err, &httpErr):
		return status == httpErr.Status
	default:
		return status == http.StatusInternalServerError
	}
}

func AsError(err error) Error {
	httpErr := Error{}
	if !errors.As(err, &httpErr) {
		return NewError(http.StatusInternalServerError, Cause(err))
	}

	return httpErr
}
