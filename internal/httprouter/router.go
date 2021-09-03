package httprouter

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

var ParamsFromContext = httprouter.ParamsFromContext

type HandlerFunc func(w http.ResponseWriter, req *http.Request) error

type Middleware func(handler HandlerFunc) HandlerFunc

type ErrorHandlerFunc func(w http.ResponseWriter, req *http.Request, err error)

type PanicHandlerFunc func(w http.ResponseWriter, req *http.Request, pv interface{})

type Router struct {
	logger       *zap.Logger
	config       config
	delegate     *httprouter.Router
	errorHandler ErrorHandlerFunc
}

func New(logger *zap.Logger, opts ...Opt) *Router {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}

	errorHandler := c.errorHandler(logger, c.verbose)

	delegate := httprouter.New()
	delegate.HandleMethodNotAllowed = true
	delegate.NotFound = adaptHandler(logger, errorHandler, func(http.ResponseWriter, *http.Request) error {
		return NewError(http.StatusNotFound)
	})
	delegate.MethodNotAllowed = adaptHandler(logger, errorHandler, func(http.ResponseWriter, *http.Request) error {
		return NewError(http.StatusMethodNotAllowed)
	})
	delegate.PanicHandler = c.panicHandler(logger, c.verbose)

	if c.globalOptions != nil {
		delegate.HandleOPTIONS = true
		delegate.GlobalOPTIONS = adaptHandler(logger, errorHandler, c.globalOptions)
	}

	return &Router{
		logger:       logger,
		config:       c,
		delegate:     delegate,
		errorHandler: errorHandler,
	}
}

func (r *Router) Handler(method, path string, handler HandlerFunc, middleware ...Middleware) {
	middleware = append(r.config.middleware, middleware...)

	for _, mw := range middleware {
		handler = mw(handler)
	}

	r.delegate.HandlerFunc(method, path, adaptHandler(r.logger, r.errorHandler, handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.delegate.ServeHTTP(w, req)
}

func adaptHandler(logger *zap.Logger, errorHandler ErrorHandlerFunc, handler HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		rw := NewResponseWriter(w)

		err := handler(rw, req)
		if err != nil {
			errorHandler(rw, req, err)
		}

		logger.Info("http request", zap.String("path", req.URL.Path), zap.Int("status", rw.StatusCode()), zap.Duration("elapsed", time.Since(start)))
	})
}
