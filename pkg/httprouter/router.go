package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var ParamsFromContext = httprouter.ParamsFromContext

type HandlerFunc func(w http.ResponseWriter, req *http.Request) error

type MiddlewareFunc func(handler HandlerFunc) HandlerFunc

type ErrorHandlerFunc func(w http.ResponseWriter, req *http.Request, verbose bool, err error)

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
		return NewError(http.StatusNotFound)
	})
	delegate.MethodNotAllowed = adaptHandler(c, func(http.ResponseWriter, *http.Request) error {
		return NewError(http.StatusMethodNotAllowed)
	})

	if c.globalOptions != nil {
		delegate.HandleOPTIONS = true
		delegate.GlobalOPTIONS = c.globalOptions
	}

	return &Router{
		config:   c,
		delegate: delegate,
	}
}

func (r *Router) Handler(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	middleware = append(r.config.middleware, middleware...)

	for _, mw := range middleware {
		handler = mw(handler)
	}

	r.delegate.HandlerFunc(method, path, adaptHandler(r.config, handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.delegate.ServeHTTP(w, req)
}

func adaptHandler(c config, handler HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)
		if err == nil {
			return
		}

		c.errorHandler(w, req, c.verbose, err)
	})
}
