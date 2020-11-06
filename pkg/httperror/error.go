package httperror

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int
	Message string
	Cause   error
}

func New(code int, opts ...ErrorOpt) Error {
	err := Error{Code: code}

	for _, opt := range opts {
		opt(&err)
	}

	if err.Message == "" {
		err.Message = http.StatusText(code)
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
