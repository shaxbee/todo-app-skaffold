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

type HandlerFunc func(w http.ResponseWriter, req *http.Request) error

type Middleware func(handler HandlerFunc) HandlerFunc

type Router struct {
	config   config
	delegate *httprouter.Router
}

func New(opts ...Opt) *Router {
	c := defaultConfig

	for _, opt := range opts {
		opt(&c)
	}

	delegate := httprouter.New()
	delegate.HandleMethodNotAllowed = true
	delegate.NotFound = adaptHandler(c, func(http.ResponseWriter, *http.Request) error {
		return httperror.New(http.StatusNotFound)
	})
	delegate.MethodNotAllowed = adaptHandler(c, func(http.ResponseWriter, *http.Request) error {
		return httperror.New(http.StatusMethodNotAllowed)
	})

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

func (r *Router) Handler(method, path string, handler HandlerFunc, middleware ...Middleware) {
	for _, mw := range middleware {
		handler = mw(handler)
	}

	r.delegate.HandlerFunc(method, path, adaptHandler(r.config, handler))
}

func adaptHandler(c config, handler HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if c.corsEnabled {
			header := w.Header()

			header.Set("Access-Control-Allow-Origin", c.corsOrigin)

			switch {
			case c.CorsRequestHeadersWildcard() && c.corsAllowCredentials:
				header.Set("Access-Control-Allow-Headers", "*, Authorization")
			case c.CorsRequestHeadersWildcard():
				header.Set("Access-Control-Allow-Headers", "*")
			default:
				header.Set("Access-Control-Allow-Headers", c.corsRequestHeaders)
			}
		}

		err := handler(w, req)
		if err == nil {
			return
		}

		httpErr := httperror.Error{}
		if !errors.As(err, &httpErr) {
			httpErr = httperror.New(http.StatusInternalServerError, httperror.Cause(err))
		}

		var debug string
		switch {
		case c.verbose && httpErr.Cause != nil:
			debug = httpErr.Cause.Error()
			log.Printf("http error %d %s: %+v", httpErr.Status, httpErr.Message, httpErr.Cause)
		case c.verbose:
			log.Printf("http error %d %s", httpErr.Status, httpErr.Message)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.Status)

		data, err := json.Marshal(api.ErrorResponse{
			Message: httpErr.Message,
			Debug:   api.PtrString(debug),
		})
		if err != nil {
			log.Printf("failed to marshal error response: %v", err)
			return
		}

		_, _ = w.Write(data)
	})
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
