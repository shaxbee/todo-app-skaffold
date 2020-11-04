package httperror

import "net/http"

type Error struct {
	Code    int
	Message string
	Cause   error
}

func New(code int, opts ...errorOpt) Error {
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

type errorOpt func(*Error)

func Message(message string) errorOpt {
	return func(e *Error) {
		e.Message = message
	}
}

func Cause(cause error) errorOpt {
	return func(e *Error) {
		e.Cause = cause
	}
}
