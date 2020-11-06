package httperror

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type (
	ErrorHandler func(http.ResponseWriter, *http.Request) error
	Middleware   func(handler ErrorHandler) http.Handler
)

func NewMiddleware(opts ...ConfigOpt) Middleware {
	c := config{}

	for _, opt := range opts {
		opt(&c)
	}

	return func(handler ErrorHandler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			err := handler(w, req)
			if err == nil {
				return
			}

			httpErr := Error{}
			if !errors.As(err, &httpErr) {
				httpErr = New(http.StatusInternalServerError, Cause(err))
			}

			var debug string
			switch {
			case c.verbose && httpErr.Cause != nil:
				debug = httpErr.Cause.Error()
				log.Printf("http error %d %s: %+v", httpErr.Code, httpErr.Message, httpErr.Cause)
			case c.verbose:
				log.Printf("http error %d %s", httpErr.Code, httpErr.Message)
			}

			data, err := json.Marshal(ErrorResponse{
				Message: httpErr.Message,
				Debug:   debug,
			})
			if err != nil {
				log.Printf("failed to marshal error response: %v", err)
			}

			w.WriteHeader(httpErr.Code)
			_, _ = w.Write(data)
		})
	}
}

type config struct {
	verbose bool
}

type ConfigOpt func(*config)

func Verbose(verbose bool) ConfigOpt {
	return func(c *config) {
		c.verbose = verbose
	}
}
