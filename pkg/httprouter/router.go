package httprouter

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/shaxbee/todo-app-skaffold/pkg/api"
	"github.com/shaxbee/todo-app-skaffold/pkg/httperror"
)

var ParamsFromContext = httprouter.ParamsFromContext

type Router struct {
	config   config
	delegate *httprouter.Router
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func New(opts ...Opt) *Router {
	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	delegate := httprouter.New()
	delegate.HandleMethodNotAllowed = true
	delegate.NotFound = notFound(c)
	delegate.MethodNotAllowed = methodNotAllowed(c)

	if c.corsEnabled {
		delegate.HandleOPTIONS = true
		delegate.GlobalOPTIONS = globalOptions(c)
	}

	return &Router{
		config:   c,
		delegate: delegate,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.delegate.ServeHTTP(w, req)
}

func (r *Router) Handler(method, path string, handler HandlerFunc) {
	r.delegate.HandlerFunc(method, path, r.adaptHandler(handler))
}

func (r *Router) adaptHandler(handler HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if r.config.corsEnabled {
			r.handleCors(w)
		}

		err := handler(w, req)
		if err == nil {
			return
		}

		httpErr := httperror.Error{}
		if !errors.As(err, &httpErr) {
			httpErr = httperror.New(http.StatusInternalServerError, httperror.Cause(err))
		}

		r.handleError(w, httpErr)
	})
}

func (r *Router) handleCors(w http.ResponseWriter) {
	header := w.Header()

	header.Set("Access-Control-Allow-Origin", r.config.corsOrigin)

	switch {
	case r.config.CorsRequestHeadersWildcard() && r.config.corsAllowCredentials:
		header.Set("Access-Control-Allow-Headers", "*, Authorization")
	case r.config.CorsRequestHeadersWildcard():
		header.Set("Access-Control-Allow-Headers", "*")
	default:
		header.Set("Access-Control-Allow-Headers", r.config.corsRequestHeaders)
	}
}

func (r *Router) handleError(w http.ResponseWriter, httpErr httperror.Error) {
	var debug string
	switch {
	case r.config.verbose && httpErr.Cause != nil:
		debug = httpErr.Cause.Error()
		log.Printf("http error %d %s: %+v", httpErr.Code, httpErr.Message, httpErr.Cause)
	case r.config.verbose:
		log.Printf("http error %d %s", httpErr.Code, httpErr.Message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)

	data, err := json.Marshal(api.ErrorResponse{
		Message: httpErr.Message,
		Debug:   debug,
	})
	if err != nil {
		log.Printf("failed to marshal error response: %v", err)
		return
	}

	_, _ = w.Write(data)
}

func notFound(c config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		data, err := json.Marshal(api.ErrorResponse{
			Message: http.StatusText(http.StatusNotFound),
		})
		if err != nil && c.verbose {
			log.Println(err)
			return
		}

		_, _ = w.Write(data)
	}
}

func methodNotAllowed(c config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		data, err := json.Marshal(api.ErrorResponse{
			Message: http.StatusText(http.StatusMethodNotAllowed),
		})
		if err != nil && c.verbose {
			log.Println(err)
			return
		}

		_, _ = w.Write(data)
	}
}

func globalOptions(c config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Access-Control-Request-Method") == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		header := w.Header()

		origin := req.Header.Get("Origin")
		requestHeaders := req.Header.Get("Access-Control-Request-Headers")

		switch {
		case c.CorsOriginWildcard() && origin != "":
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set("Vary", origin)
		case c.CorsOriginWildcard():
			header.Set("Access-Control-Allow-Origin", c.corsOrigin)
		default:
			return
		}

		switch {
		case c.CorsRequestHeadersWildcard() && requestHeaders != "":
			header.Set("Access-Control-Allow-Headers", requestHeaders)
		case requestHeaders != "":
			header.Set("Access-Control-Allow-Headers", c.corsRequestHeaders)
		}

		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))

		if c.corsAllowCredentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}

		if c.corsMaxAge != "" {
			header.Set("Access-Control-Max-Age", c.corsMaxAge)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
